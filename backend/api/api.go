package api

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats.go"

	"gitlab.com/timeterm/timeterm/nats-manager/sdk"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/mq"
	"gitlab.com/timeterm/timeterm/backend/secrets"
	"gitlab.com/timeterm/timeterm/backend/templates"
)

type Server struct {
	db   *database.Wrapper
	log  logr.Logger
	echo *echo.Echo
	mqw  *mq.Wrapper
	secr *secrets.Wrapper
	nm   *nmsdk.Client
}

func newEcho(log logr.Logger) (*echo.Echo, error) {
	var err error

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = newEchoLogrLogger(log)
	e.Use(middleware.CORS())

	e.Renderer, err = templates.Load()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewServer(db *database.Wrapper, log logr.Logger, nc *nats.Conn, secr *secrets.Wrapper) (Server, error) {
	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:   db,
		log:  log,
		echo: e,
		secr: secr,
		mqw:  mq.NewWrapper(nc),
		nm:   nmsdk.NewClient(nc),
	}
	server.registerRoutes()

	authnr, err := authn.New(db, log)
	if err != nil {
		return server, err
	}
	authnr.RegisterRoutes(e.Group(""))

	return server, nil
}

func (s *Server) registerRoutes() {
	g := s.echo.Group("")
	g.Use(authn.UserLoginMiddleware(s.db, s.log))

	userGroup := g.Group("/user")
	userGroup.GET("/me", s.getCurrentUser)
	userGroup.PATCH("/:id", s.patchUser)

	devGroup := g.Group("/device")
	devGroup.GET("", s.getDevices)
	devGroup.DELETE("", s.deleteDevices)
	devGroup.POST("/restart", s.rebootDevices)
	devGroup.GET("/:id", s.getDevice)
	devGroup.POST("/:id/restart", s.rebootDevice)
	devGroup.PATCH("/:id", s.patchDevice)
	devGroup.DELETE("/:id", s.deleteDevice)
	devGroup.GET("/registrationconfig", s.getRegistrationConfig)

	registrationLoginMiddleware := authn.DeviceRegistrationLoginMiddleware(s.db, s.log)
	s.echo.POST("/device", registrationLoginMiddleware(s.createDevice))

	devConfigGroup := s.echo.Group("/device/:id/config")
	devConfigGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log))
	devConfigGroup.GET("/natscreds", s.generateNATSCredentials)

	devHeartbeatGroup := s.echo.Group("/device/:id/heartbeat")
	devHeartbeatGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log))
	devHeartbeatGroup.POST("", s.updateLastHeartbeat)

	orgGroup := g.Group("/organization")
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	stdGroup := g.Group("/student")
	stdGroup.GET("/:id", s.getStudent)
	stdGroup.PATCH("/:id", s.patchStudent)
	stdGroup.GET("", s.getStudents)
	stdGroup.POST("", s.createStudent)
	stdGroup.DELETE("", s.deleteStudents)

	netServGroup := g.Group("/networking/service")
	netServGroup.GET("", s.getNetworkingServices)
	netServGroup.GET("/:id", s.getNetworkingService)
	netServGroup.POST("", s.createNetworkingService)
	netServGroup.PUT("/:id", s.replaceNetworkingService)
	netServGroup.DELETE("/:id", s.deleteNetworkingService)
}

func (s *Server) Run(ctx context.Context) error {
	const shutdownTimeout = time.Second * 30

	errc := make(chan error)
	go func() {
		const serveAddr = ":1323"
		s.log.Info("serving", "addr", serveAddr)
		errc <- s.echo.Start(serveAddr)
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		s.log.Info("shutting down API server", "timeout", shutdownTimeout)
		if err := s.echo.Shutdown(shutdownCtx); err != nil {
			s.log.Error(err, "failed to gracefully shut down API server")
		}
		s.log.Info("done shutting down API server")

		return ctx.Err()
	}
}

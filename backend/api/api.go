package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats.go"

	nmsdk "gitlab.com/timeterm/timeterm/nats-manager/pkg/sdk"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/messages"
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
	msgw *messages.Wrapper
}

func newEcho(log logr.Logger) (*echo.Echo, error) {
	var err error

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = newEchoLogrLogger(log)
	e.Use(middleware.Recover(), middleware.CORS())

	e.Renderer, err = templates.Load()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewServer(log logr.Logger, db *database.Wrapper, secr *secrets.Wrapper) (Server, error) {
	log = log.WithName("Server")

	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	acr, err := nmsdk.NewAppCredsRetrieverFromEnv("backend")
	if err != nil {
		return Server{}, fmt.Errorf("could not create (NATS) app credentials retriever: %w", err)
	}

	nc, err := nats.Connect(os.Getenv("NATS_URL"),
		nats.UserJWT(acr.NatsCredsCBs()),
		// Never stop trying to reconnect.
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return Server{}, fmt.Errorf("could not connect to NATS: %w", err)
	}

	mqw, err := mq.NewWrapper(log, db)
	if err != nil {
		return Server{}, fmt.Errorf("could not create NATS wrapper: %w", err)
	}

	server := Server{
		db:   db,
		log:  log,
		echo: e,
		secr: secr,
		mqw:  mqw,
		nm:   nmsdk.NewClient(nc),
		msgw: messages.NewWrapper(log, db, secr),
	}
	server.registerRoutes()

	authnr, err := authn.New(log, db, secr)
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
	devGroup.GET("/:id", s.getDevice)
	devGroup.PATCH("/:id", s.patchDevice)
	devGroup.DELETE("/:id", s.deleteDevice)
	devGroup.POST("/:id/restart", s.rebootDevice)
	devGroup.GET("/registrationconfig", s.getRegistrationConfig)
	devGroup.POST("/restart", s.rebootDevices)

	registrationLoginMiddleware := authn.DeviceRegistrationLoginMiddleware(s.db, s.log)
	s.echo.POST("/device", registrationLoginMiddleware(s.createDevice))

	devConfigGroup := s.echo.Group("/device/:id/config")
	devConfigGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log))
	devConfigGroup.GET("/natscreds", s.generateNATSCredentials)
	devConfigGroup.GET("/networks", s.getAllNetworkingServices)

	devHeartbeatGroup := s.echo.Group("/device/:id/heartbeat")
	devHeartbeatGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log))
	devHeartbeatGroup.PUT("", s.updateLastHeartbeat)

	msgGroup := s.echo.Group("/message")
	msgGroup.GET("", s.getAdminMessages)
	msgGroup.GET("/:sec/:nanosec", s.getAdminMessage)

	orgGroup := g.Group("/organization")
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	stdGroup := g.Group("/student")
	stdGroup.GET("", s.getStudents)
	stdGroup.POST("", s.createStudent)
	stdGroup.DELETE("", s.deleteStudents)
	stdGroup.GET("/:id", s.getStudent)
	stdGroup.PATCH("/:id", s.patchStudent)

	netServGroup := g.Group("/networking/service")
	netServGroup.GET("", s.getNetworkingServices)
	netServGroup.POST("", s.createNetworkingService)
	netServGroup.GET("/:id", s.getNetworkingService)
	netServGroup.PUT("/:id", s.replaceNetworkingService)
	netServGroup.DELETE("/:id", s.deleteNetworkingService)

	zappGroup := s.echo.Group("/zermelo/appointment")
	zappGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log), authn.StudentLoginMiddleware(s.db, s.log))
	zappGroup.GET("", s.getZermeloAppointments)

	zenrGroup := s.echo.Group("/zermelo/enrollment")
	zenrGroup.Use(authn.DeviceLoginMiddleware(s.db, s.log), authn.StudentLoginMiddleware(s.db, s.log))
	zenrGroup.POST("", s.enrollZermelo)

	zconnGroup := g.Group("/zermelo/connect")
	zconnGroup.POST("", s.connectZermeloOrganization)
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

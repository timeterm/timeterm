package api

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/labstack/echo"

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
}

func newEcho(log logr.Logger) (*echo.Echo, error) {
	var err error

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = newEchoLogrLogger(log)

	e.Renderer, err = templates.Load()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewServer(db *database.Wrapper, log logr.Logger, mqw *mq.Wrapper, secr *secrets.Wrapper) (Server, error) {
	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:   db,
		log:  log,
		echo: e,
		mqw:  mqw,
		secr: secr,
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
	g.Use(authn.Middleware(s.db, s.log))

	userGroup := s.echo.Group("/user")
	userGroup.GET("/me", s.getCurrentUser)
	userGroup.PATCH("/:id", s.patchUser)

	devGroup := s.echo.Group("/device")
	devGroup.GET("/", s.getDevices)
	devGroup.DELETE("/", s.deleteDevices)
	devGroup.POST("/", s.createDevice)
	devGroup.POST("/restart", s.rebootDevices)
	devGroup.GET("/:id", s.getDevice)
	devGroup.POST("/:id/restart", s.rebootDevice)
	devGroup.PATCH("/:id", s.patchDevice)
	devGroup.DELETE("/:id", s.deleteDevice)

	orgGroup := s.echo.Group("/organization")
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	stdGroup := s.echo.Group("/student")
	stdGroup.GET("/:id", s.getStudent)
	stdGroup.PATCH("/:id", s.patchStudent)
	stdGroup.GET("/", s.getStudents)
	stdGroup.POST("/", s.createStudent)
	stdGroup.DELETE("/", s.deleteStudents)

	ethServGroup := s.echo.Group("/ethernet/service")
	ethServGroup.GET("/:id", s.getEthernetService)
}

func (s *Server) Run(ctx context.Context) error {
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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		_ = s.echo.Shutdown(shutdownCtx)
		return ctx.Err()
	}
}

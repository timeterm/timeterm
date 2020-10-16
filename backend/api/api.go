package api

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/broker"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/templates"
)

type Server struct {
	db       *database.Wrapper
	log      logr.Logger
	echo     *echo.Echo
	apiGroup *echo.Group
	brw      *broker.Wrapper
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

func NewServer(db *database.Wrapper, log logr.Logger, brw *broker.Wrapper) (Server, error) {
	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:       db,
		log:      log,
		echo:     e,
		apiGroup: e.Group("/api"),
		brw:      brw,
	}
	server.registerRoutes()

	authnr, err := authn.New(db, log)
	if err != nil {
		return server, err
	}
	authnr.RegisterRoutes(server.apiGroup)

	return server, nil
}

func (s *Server) registerRoutes() {
	g := s.apiGroup.Group("")
	g.Use(authn.Middleware(s.db, s.log))

	g.GET("/user/me", s.getCurrentUser)
	g.PATCH("/user/:id", s.patchUser)

	g.GET("/device", s.getDevices)
	g.DELETE("/device", s.deleteDevices)
	g.POST("/device", s.createDevice)
	g.POST("/device/restart", s.rebootDevices)
	g.GET("/device/:id", s.getDevice)
	g.POST("/device/:id/restart", s.rebootDevice)
	g.PATCH("/device/:id", s.patchDevice)
	g.DELETE("/device/:id", s.deleteDevice)

	orgGroup := s.apiGroup.Group("/organization")
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	g.GET("/student/:id", s.getStudent)
	g.PATCH("/student/:id", s.patchStudent)
	g.GET("/student", s.getStudents)
	g.POST("/student", s.createStudent)
	g.DELETE("/student", s.deleteStudents)
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

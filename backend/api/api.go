package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gitlab.com/timeterm/timeterm/backend/database"
)

type Server struct {
	db   *database.Wrapper
	echo *echo.Echo
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func NewServer(db *database.Wrapper) Server {
	server := Server{
		db:   db,
		echo: newEcho(),
	}
	server.registerRoutes()

	return server
}

func (s *Server) registerRoutes() {
	s.echo.GET("/organization/:id", s.getOrganization)
	s.echo.GET("/student/:id", s.getStudent)
	s.echo.GET("/device/:id", s.getDevice)
}

func (s *Server) getOrganization(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	dbOrg, err := s.db.GetOrganization(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read organization from database")
	}

	apiOrg := OrganizationFrom(dbOrg)
	return c.JSON(http.StatusOK, apiOrg)
}

func (s *Server) getStudent(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	dbStudent, err := s.db.GetStudent(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read student from database")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

func (s *Server) getDevice(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	dbDevice, err := s.db.GetDevice(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read device from database")
	}

	apiDevice := DeviceFrom(dbDevice)
	return c.JSON(http.StatusOK, apiDevice)
}

func (s *Server) Run(ctx context.Context) error {
	errc := make(chan error)
	go func() {
		errc <- s.echo.Start(":1323")
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		_ = s.echo.Close()
		return ctx.Err()
	}
}

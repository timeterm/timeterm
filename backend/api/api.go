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
	s.echo.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Leutige boel")
	})

	s.echo.GET("/user/:id", func(c echo.Context) error {
		id := c.Param("id")

		uid, err := uuid.Parse(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}

		dbUser, err := s.db.ReadUser(c.Request().Context(), uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not read user")
		}

		apiUser := UserFrom(dbUser)
		return c.JSON(http.StatusOK, apiUser)
	})
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

package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/templates"
)

type Server struct {
	db   *database.Wrapper
	log  logr.Logger
	echo *echo.Echo
}

func newEcho() (*echo.Echo, error) {
	var err error

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Renderer, err = templates.Load()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewServer(db *database.Wrapper, log logr.Logger) (Server, error) {
	e, err := newEcho()
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:   db,
		log:  log,
		echo: e,
	}
	server.registerRoutes()

	authnr, err := authn.New(db, log)
	if err != nil {
		return server, err
	}
	authnr.RegisterRoutes(server.echo)

	return server, nil
}

func (s *Server) registerRoutes() {
	s.echo.GET("/device/:id", s.getDevice)
	s.echo.GET("/device", s.getDevices)

	orgGroup := s.echo.Group("/organization")
	orgGroup.POST("/:organization/student", s.createStudent)
	orgGroup.POST("/:organization/device", s.createDevice)
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	s.echo.GET("/student/:id", s.getStudent)
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

func (s *Server) getDevices(c echo.Context) error {
	dbDevices, err := s.db.GetDevices(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read devices from database")
	}

	apiDevices := DevicesFrom(dbDevices)

	return c.JSON(http.StatusOK, apiDevices)
}

func (s *Server) createStudent(c echo.Context) error {
	organizationID := c.Param("organization")

	uid, err := uuid.Parse(organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	dbStudent, err := s.db.CreateStudent(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not create student")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create student")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

func (s *Server) createDevice(c echo.Context) error {
	organizationID := c.Param("organization")

	uid, err := uuid.Parse(organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var dev Device
	err = c.Bind(&dev)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not bind data")
	}

	dbDevice, err := s.db.CreateDevice(c.Request().Context(), uid, dev.Name, dev.Status)
	if err != nil {
		s.log.Error(err, "could not create device")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create device")
	}

	apiDevice := DeviceFrom(dbDevice)
	return c.JSON(http.StatusOK, apiDevice)
}

func (s *Server) patchOrganization(c echo.Context) error {
	organizationID := c.Param("id")

	uid, err := uuid.Parse(organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	bytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read request body")
	}

	oldDBOrganization, err := s.db.GetOrganization(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not read organization from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read organization from database")
	}

	oldAPIOrganization := OrganizationFrom(oldDBOrganization)

	jsonOrganization, err := json.Marshal(oldAPIOrganization)
	if err != nil {
		s.log.Error(err, "could not marshal the old organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not marshal the old organization")
	}

	newJSONOrganization, err := jsonpatch.MergePatch(bytes, jsonOrganization)
	if err != nil {
		s.log.Error(err, "could not patch the organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not patch the organization")
	}

	var newAPIOrganization Organization
	err = json.Unmarshal(newJSONOrganization, &newAPIOrganization)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not unmarshal patched organization")
	}

	newAPIOrganization.ID = uid

	newDBOrganization := OrganisationToDB(newAPIOrganization)

	err = s.db.ReplaceOrganization(c.Request().Context(), newDBOrganization)
	if err != nil {
		s.log.Error(err, "could not update the organization in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update the organization in the database")
	}

	return c.JSON(http.StatusOK, newAPIOrganization)
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

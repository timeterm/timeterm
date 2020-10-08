package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/templates"
)

type Server struct {
	db       *database.Wrapper
	log      logr.Logger
	echo     *echo.Echo
	apiGroup *echo.Group
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

func NewServer(db *database.Wrapper, log logr.Logger) (Server, error) {
	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:       db,
		log:      log,
		echo:     e,
		apiGroup: e.Group("/api"),
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
	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-Api-Key",
		AuthScheme: "",
		Validator: func(key string, c echo.Context) (bool, error) {
			token, err := uuid.Parse(key)
			if err != nil {
				return false, echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
			}

			user, err := s.db.GetUserByToken(c.Request().Context(), token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}

				s.log.Error(err, "failed to get user by token")
				return false, echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
			}

			c.Set("user", user)

			return true, nil
		},
	}))

	g.GET("/user/me", s.getCurrentUser)

	g.GET("/device/:id", s.getDevice)
	g.DELETE("/device/:id", s.deleteDevice)
	g.GET("/device", s.getDevices)

	orgGroup := s.apiGroup.Group("/organization")
	orgGroup.POST("/:organization/student", s.createStudent)
	orgGroup.POST("/:organization/device", s.createDevice)
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	g.GET("/student/:id", s.getStudent)
}

func (s *Server) getCurrentUser(c echo.Context) error {
	dbUser := c.Get("user").(database.User)

	apiUser := UserFrom(dbUser)
	return c.JSON(http.StatusOK, apiUser)
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbDevice, err := s.db.GetDevice(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read device from database")
	}

	apiDevice := DeviceFrom(dbDevice)
	return c.JSON(http.StatusOK, apiDevice)
}

func (s *Server) deleteDevice(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dev, err := s.db.GetDevice(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not get device")
		return echo.NewHTTPError(http.StatusBadRequest, "Could not get device")
	}

	user := c.Get("user").(database.User)
	if dev.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Device does not belong to user's organization")
	}

	err = s.db.DeleteDevice(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not delete device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete device")
	}

	return c.NoContent(http.StatusNoContent)
}

type getDevicesParams struct {
	Offset     *uint64 `query:"offset"`
	MaxAmount  *uint64 `query:"maxAmount"`
	SearchName *string `query:"searchName"`
}

func (s *Server) getDevices(c echo.Context) error {
	var params getDevicesParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user := c.Get("user").(database.User)

	dbDevices, err := s.db.GetDevices(c.Request().Context(), database.GetDevicesOpts{
		OrganizationID: user.OrganizationID,
		Limit:          params.MaxAmount,
		Offset:         params.Offset,
		NameSearch:     params.SearchName,
	})
	if err != nil {
		s.log.Error(err, "could not get devices")
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read devices from database")
	}

	apiDevices := PaginatedDevicesFrom(dbDevices)

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

	dbDevice, err := s.db.CreateDevice(c.Request().Context(), uid, dev.Name, database.DeviceStatusOffline)
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
		const serveAddr = ":1323"
		s.log.Info("serving", "addr", serveAddr)
		errc <- s.echo.Start(serveAddr)
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		_ = s.echo.Close()
		return ctx.Err()
	}
}

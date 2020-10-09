package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/nats-io/nats.go"
	timetermpb "gitlab.com/timeterm/timeterm/proto/go"
	"google.golang.org/protobuf/proto"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/templates"
)

type Server struct {
	db       *database.Wrapper
	log      logr.Logger
	echo     *echo.Echo
	apiGroup *echo.Group
	enc      *nats.EncodedConn
}

type protoEncoder struct{}

func (p protoEncoder) Encode(_ string, v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("v is not proto.Message")
	}
	return proto.Marshal(msg)
}

func (p protoEncoder) Decode(_ string, data []byte, vPtr interface{}) error {
	msg, ok := vPtr.(proto.Message)
	if !ok {
		return errors.New("vPtr is not proto.Message")
	}
	return proto.Unmarshal(data, msg)
}

func init() {
	nats.RegisterEncoder("proto", &protoEncoder{})
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

func NewServer(db *database.Wrapper, log logr.Logger, nc *nats.Conn) (Server, error) {
	e, err := newEcho(log)
	if err != nil {
		return Server{}, err
	}

	enc, err := nats.NewEncodedConn(nc, "proto")
	if err != nil {
		return Server{}, err
	}

	server := Server{
		db:       db,
		log:      log,
		echo:     e,
		apiGroup: e.Group("/api"),
		enc:      enc,
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

	g.GET("/device", s.getDevices)
	g.GET("/device/:id", s.getDevice)
	g.POST("/device/:id/restart", s.rebootDevice)
	g.DELETE("/device/:id", s.deleteDevice)

	orgGroup := s.apiGroup.Group("/organization")
	orgGroup.POST("/:organization/student", s.createStudent)
	orgGroup.POST("/:organization/device", s.createDevice)
	orgGroup.PATCH("/:id", s.patchOrganization)
	orgGroup.GET("/:id", s.getOrganization)

	g.GET("/student/:id", s.getStudent)
}

func (s *Server) getCurrentUser(c echo.Context) error {
	dbUser, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

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

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}
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

func (s *Server) rebootDevice(c echo.Context) error {
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

	err = s.enc.Publish(fmt.Sprintf("FEDEV.%s.REBOOT", id), &timetermpb.RebootMessage{
		DeviceId: id,
	})
	if err != nil {
		s.log.Error(err, "could not publish reboot message")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not send reboot message")
	}

	return c.NoContent(http.StatusOK)
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

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

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

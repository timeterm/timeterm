package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbOrg, err := s.db.GetOrganization(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read organization from database")
	}

	apiOrg := OrganizationFrom(dbOrg)
	return c.JSON(http.StatusOK, apiOrg)
}

func (s *Server) getStudent(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbStudent, err := s.db.GetStudent(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read student from database")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

type getStudentsParams struct {
	Offset    *uint64 `query:"offset"`
	MaxAmount *uint64 `query:"maxAmount"`
}

func (s *Server) getStudents(c echo.Context) error {
	var params getStudentsParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	dbStudents, err := s.db.GetStudents(c.Request().Context(), database.GetStudentsOpts{
		OrganizationID: user.OrganizationID,
		Limit:          params.MaxAmount,
		Offset:         params.Offset,
	})
	if err != nil {
		s.log.Error(err, "could not read students from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read students from database")
	}

	apiStudents := PaginatedStudentsFrom(dbStudents)
	return c.JSON(http.StatusOK, apiStudents)
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

type deleteDevicesParams struct {
	DeviceIDs []uuid.UUID `json:"deviceIds"`
}

func (s *Server) deleteDevices(c echo.Context) error {
	var params deleteDevicesParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	allInOrg, err := s.db.AreDevicesInOrganization(c.Request().Context(), user.OrganizationID, params.DeviceIDs...)
	if err != nil {
		s.log.Error(err, "could not get devices in organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve device information")
	}
	if !allInOrg {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not all devices are in user's organization")
	}

	err = s.db.DeleteDevices(c.Request().Context(), params.DeviceIDs)
	if err != nil {
		s.log.Error(err, "could not delete devices")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete devices")
	}

	return c.NoContent(http.StatusNoContent)
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

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	if dev.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Device does not belong to user's organization")
	}

	err = s.brw.RebootDevice(uid)
	if err != nil {
		s.log.Error(err, "could not publish reboot message")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not send reboot message")
	}

	return c.NoContent(http.StatusOK)
}

type rebootDevicesParams struct {
	DeviceIDs []uuid.UUID `json:"deviceIds"`
}

func (s *Server) rebootDevices(c echo.Context) error {
	var params rebootDevicesParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	allInOrg, err := s.db.AreDevicesInOrganization(c.Request().Context(), user.OrganizationID, params.DeviceIDs...)
	if err != nil {
		s.log.Error(err, "could not get devices in organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve device information")
	}
	if !allInOrg {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not all devices are in user's organization")
	}

	for _, deviceID := range params.DeviceIDs {
		err = s.brw.RebootDevice(deviceID)
		if err != nil {
			s.log.Error(err, "could not publish reboot message")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not send reboot message")
		}
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read devices from database")
	}

	apiDevices := PaginatedDevicesFrom(dbDevices)

	return c.JSON(http.StatusOK, apiDevices)
}

func (s *Server) createStudent(c echo.Context) error {
	organizationID := c.Param("organization")

	uid, err := uuid.Parse(organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbStudent, err := s.db.CreateStudent(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not create student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create student")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

func (s *Server) patchDevice(c echo.Context) error {
	deviceID := c.Param("id")

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	uid, err := uuid.Parse(deviceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	patchData, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	oldDBDevice, err := s.db.GetDevice(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not read device from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read device from database")
	}

	if oldDBDevice.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Device not in organization")
	}

	oldAPIDevice := DeviceFrom(oldDBDevice)

	jsonDevice, err := json.Marshal(oldAPIDevice)
	if err != nil {
		s.log.Error(err, "could not marshal the old device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not marshal the old device")
	}

	newJSONDevice, err := jsonpatch.MergePatch(jsonDevice, patchData)
	if err != nil {
		s.log.Error(err, "could not patch the device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch the device")
	}

	var newAPIDevice Device
	err = json.Unmarshal(newJSONDevice, &newAPIDevice)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched device")
	}

	newAPIDevice.ID = oldDBDevice.ID
	newAPIDevice.PrimaryStatus = oldAPIDevice.PrimaryStatus
	newAPIDevice.SecondaryStatus = oldAPIDevice.SecondaryStatus
	newAPIDevice.OrganizationID = oldDBDevice.OrganizationID

	newDBDevice := DeviceToDB(newAPIDevice)
	newDBDevice.LastHeartbeat = oldDBDevice.LastHeartbeat

	err = s.db.ReplaceDevice(c.Request().Context(), newDBDevice)
	if err != nil {
		s.log.Error(err, "could not update the device in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update the device in the database")
	}

	return c.JSON(http.StatusOK, newAPIDevice)
}

func (s *Server) patchStudent(c echo.Context) error {
	studentID := c.Param("id")

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	uid, err := uuid.Parse(studentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	patchData, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	oldDBStudent, err := s.db.GetStudent(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not read student from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read student from database")
	}

	if oldDBStudent.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not in organization")
	}

	oldAPIStudent := StudentFrom(oldDBStudent)

	jsonStudent, err := json.Marshal(oldAPIStudent)
	if err != nil {
		s.log.Error(err, "could not marshal the old student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not marshal the old student")
	}

	newJSONStudent, err := jsonpatch.MergePatch(jsonStudent, patchData)
	if err != nil {
		s.log.Error(err, "could not patch the student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch the student")
	}

	var newAPIStudent Student
	err = json.Unmarshal(newJSONStudent, &newAPIStudent)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched student")
	}

	newAPIStudent.ID = oldDBStudent.ID
	newAPIStudent.OrganizationID = oldDBStudent.OrganizationID

	newDBStudent := StudentToDB(newAPIStudent)

	err = s.db.ReplaceStudent(c.Request().Context(), newDBStudent)
	if err != nil {
		s.log.Error(err, "could not update the student in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update the student in the database")
	}

	return c.JSON(http.StatusOK, newAPIStudent)
}

func (s *Server) createDevice(c echo.Context) error {
	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Not authenticated")
	}

	var dev Device
	err := c.Bind(&dev)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not bind data")
	}

	dbDevice, err := s.db.CreateDevice(c.Request().Context(),
		user.OrganizationID, dev.Name, database.DeviceStatusNotActivated,
	)
	if err != nil {
		s.log.Error(err, "could not create device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create device")
	}

	apiDevice := DeviceFrom(dbDevice)
	return c.JSON(http.StatusOK, apiDevice)
}

func (s *Server) patchOrganization(c echo.Context) error {
	organizationID := c.Param("id")

	uid, err := uuid.Parse(organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	patchData, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	oldDBOrganization, err := s.db.GetOrganization(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not read organization from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read organization from database")
	}

	oldAPIOrganization := OrganizationFrom(oldDBOrganization)

	jsonOrganization, err := json.Marshal(oldAPIOrganization)
	if err != nil {
		s.log.Error(err, "could not marshal the old organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not marshal the old organization")
	}

	newJSONOrganization, err := jsonpatch.MergePatch(jsonOrganization, patchData)
	if err != nil {
		s.log.Error(err, "could not patch the organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch the organization")
	}

	var newAPIOrganization Organization
	err = json.Unmarshal(newJSONOrganization, &newAPIOrganization)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched organization")
	}

	newAPIOrganization.ID = oldDBOrganization.ID
	newDBOrganization := OrganisationToDB(newAPIOrganization)

	err = s.db.ReplaceOrganization(c.Request().Context(), newDBOrganization)
	if err != nil {
		s.log.Error(err, "could not update the organization in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update the organization in the database")
	}

	return c.JSON(http.StatusOK, newAPIOrganization)
}

type deleteStudentsParams struct {
	StudentDs []uuid.UUID `json:"studentIds"`
}

func (s *Server) deleteStudents(c echo.Context) error {
	var params deleteStudentsParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	allInOrg, err := s.db.AreStudentsInOrganization(c.Request().Context(), user.OrganizationID, params.StudentDs...)
	if err != nil {
		s.log.Error(err, "could not get students in organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve student information")
	}
	if !allInOrg {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not all devices are in user's organization")
	}

	err = s.db.DeleteStudents(c.Request().Context(), params.StudentDs)
	if err != nil {
		s.log.Error(err, "could not delete students")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete students")
	}

	return c.NoContent(http.StatusNoContent)
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

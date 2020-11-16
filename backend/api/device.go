package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	jsonpatch "github.com/evanphx/json-patch/v5"
	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
)

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

	err = s.mqw.RebootDevice(uid)
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
		err = s.mqw.RebootDevice(deviceID)
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

// TODO(rutgerbrf): should be authenticated using a device registration token, not a user token.
func (s *Server) createDevice(c echo.Context) error {
	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	var dev Device
	err := c.Bind(&dev)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not bind data")
	}

	dbDevice, token, err := s.db.CreateDevice(c.Request().Context(),
		user.OrganizationID, dev.Name, database.DeviceStatusNotActivated,
	)
	if err != nil {
		s.log.Error(err, "could not create device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create device")
	}

	rsp := CreateDeviceResponseFrom(dbDevice, token)
	return c.JSON(http.StatusOK, rsp)
}

func (s *Server) getNATSCredentials(c echo.Context) error {
	_, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	// TODO(rutgerbrf): use Timeterm nats-manager to retrieve credentials

	return nil
}

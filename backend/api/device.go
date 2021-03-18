package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	dbDevice, err := s.db.GetDevice(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read device from database")
	}

	if user.OrganizationID != dbDevice.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Device does not belong to user's organization")
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
	paginationParams
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

func (s *Server) createDevice(c echo.Context) error {
	org, ok := authn.OrganizationFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	var dev Device
	err := c.Bind(&dev)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not bind data")
	}

	dbDevice, token, err := s.db.CreateDevice(c.Request().Context(), org.ID, dev.Name)
	if err != nil {
		s.log.Error(err, "could not create device")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create device")
	}

	err = s.nm.ProvisionNewDevice(c.Request().Context(), dbDevice.ID)
	if err != nil {
		s.log.Error(err, "could not provision new device")
		if err = s.db.DeleteDevice(c.Request().Context(), dbDevice.ID); err != nil {
			s.log.Error(err, "could not delete device after failed provisioning attempt")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not provision new device")
	}

	rsp := CreateDeviceResponseFrom(dbDevice, token)
	return c.JSON(http.StatusOK, rsp)
}

func (s *Server) generateNATSCredentials(c echo.Context) error {
	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	start := time.Now()
	creds, err := s.nm.GenerateDeviceCredentials(c.Request().Context(), dev.ID)
	if err != nil {
		s.log.Error(err, "could not generate NATS credentials")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not generate NATS credentials")
	}
	s.log.V(1).Info("generated NATS credentials", "took", time.Since(start))

	rsp := GenerateNATSCredentialsResponse{
		Credentials: creds,
	}
	return c.JSON(http.StatusOK, rsp)
}

func (s *Server) getRegistrationConfig(c echo.Context) error {
	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	token, err := s.db.CreateDeviceRegistrationToken(c.Request().Context(), user.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create token")
	}

	apiNetworkingServices, err := s.apiGetAllNetworkingServices(c.Request().Context(), user.OrganizationID)
	if err != nil {
		return err
	}

	rsp := RegistrationConfig{
		Token:              token,
		OrganizationID:     user.OrganizationID,
		NetworkingServices: apiNetworkingServices,
	}
	return c.JSON(http.StatusOK, rsp)
}

func (s *Server) updateLastHeartbeat(c echo.Context) error {
	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	deviceID := c.Param("id")
	uid, err := uuid.Parse(deviceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	if dev.ID != uid {
		return echo.NewHTTPError(http.StatusBadRequest, "ID mismatch")
	}

	err = s.db.ReplaceDeviceHeartbeat(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not update last heartbeat")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update last heartbeat")
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) getAllNetworkingServices(c echo.Context) error {
	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	deviceID := c.Param("id")
	uid, err := uuid.Parse(deviceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	if dev.ID != uid {
		return echo.NewHTTPError(http.StatusBadRequest, "ID mismatch")
	}

	apiNetworkingServices, err := s.apiGetAllNetworkingServices(c.Request().Context(), dev.OrganizationID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, apiNetworkingServices)
}

func (s *Server) apiGetAllNetworkingServices(ctx context.Context, organizationID uuid.UUID) ([]NetworkingService, error) {
	dbNetworkingServices, err := s.db.GetAllNetworkingServices(ctx, organizationID)
	if err != nil {
		s.log.Error(err, "could not read networking services from database")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Could not read networking services from database")
	}

	apiNetworkingServices := make([]NetworkingService, len(dbNetworkingServices))

	for i, networkingService := range dbNetworkingServices {
		uid := networkingService.ID
		secretNetworkingService, err := s.secr.GetNetworkingService(uid)
		if err != nil {
			s.log.Error(err, "could not read secret networking service")
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Could not read secret networking service")
		}

		apiNetworkingService := NetworkingServiceFrom(secretNetworkingService, networkingService)

		apiNetworkingServices[i] = apiNetworkingService
	}

	return apiNetworkingServices, nil
}

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
)

type getNetworkingServicesParams struct {
	paginationParams
}

func (s *Server) getNetworkingServices(c echo.Context) error {
	var params getNetworkingServicesParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	dbNetworkingServices, err := s.db.GetNetworkingServices(c.Request().Context(), database.GetNetworkingServicesOpts{
		OrganizationID: user.OrganizationID,
		Limit:          params.MaxAmount,
		Offset:         params.Offset,
	})
	if err != nil {
		s.log.Error(err, "could not get networking services from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read networking services from database")
	}

	apiNetworkingServices := make([]NetworkingService, len(dbNetworkingServices.NetworkingServices))

	for i, networkingService := range dbNetworkingServices.NetworkingServices {
		uid := networkingService.ID
		secretNetworkingService, err := s.secr.GetNetworkingService(uid)
		if err != nil {
			s.log.Error(err, "could not read secret networking service")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not read secret networking service")
		}

		apiNetworkingService := NetworkingServiceFrom(secretNetworkingService, networkingService)

		apiNetworkingServices[i] = apiNetworkingService
	}

	apiPaginatedNetworkingServices := PaginatedNetworkingServicesFrom(dbNetworkingServices, apiNetworkingServices)

	return c.JSON(http.StatusOK, apiPaginatedNetworkingServices)
}

func (s *Server) getNetworkingService(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbNetworkingService, err := s.db.GetNetworkingService(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not get networking service")
		return echo.NewHTTPError(http.StatusBadRequest, "Could not get networking service")
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	if dbNetworkingService.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Networking service does not belong to user's organization")
	}

	secretNetworkingService, err := s.secr.GetNetworkingService(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read secret networking service")
	}

	apiNetworkingService := NetworkingServiceFrom(secretNetworkingService, dbNetworkingService)
	return c.JSON(http.StatusOK, apiNetworkingService)
}

func (s *Server) replaceNetworkingService(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	reqBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	dbNetworkingService, err := s.db.GetNetworkingService(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not get networking service")
		return echo.NewHTTPError(http.StatusBadRequest, "Could not get networking service")
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	if dbNetworkingService.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Networking service does not belong to user's organization")
	}

	var oldNetworkingService NetworkingService

	err = json.Unmarshal(reqBody, &oldNetworkingService)
	if err != nil {
		s.log.Error(err, "could not unmarshal request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal request body")
	}

	if dbNetworkingService.Name != oldNetworkingService.Name {
		if err = s.db.ReplaceNetworkingService(c.Request().Context(), database.NetworkingService{
			ID:             oldNetworkingService.ID,
			OrganizationID: oldNetworkingService.OrganizationID,
			Name:           oldNetworkingService.Name,
		}); err != nil {
			s.log.Error(err, "could not update networking service")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not update networking service")
		}
	}

	oldProtoNetworkingService := NetworkingServiceToProto(oldNetworkingService)

	err = s.secr.UpsertNetworkingService(uid, oldProtoNetworkingService)
	if err != nil {
		s.log.Error(err, "could not update secret networking service")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update secret networking service")
	}

	s.mqw.NetworkingConfigUpdated(user.OrganizationID)

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) createNetworkingService(c echo.Context) error {
	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Not authenticated")
	}

	var ns NetworkingService
	err := c.Bind(&ns)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not bind data")
	}

	dbNetworkingService, err := s.db.CreateNetworkingService(c.Request().Context(),
		user.OrganizationID, ns.Name,
	)
	if err != nil {
		s.log.Error(err, "could not create database networking service")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create database networking service")
	}

	ns.ID = dbNetworkingService.ID
	ns.OrganizationID = dbNetworkingService.OrganizationID
	ns.Name = dbNetworkingService.Name

	secretNS := NetworkingServiceToProto(ns)

	err = s.secr.UpsertNetworkingService(ns.ID, secretNS)
	if err != nil {
		s.log.Error(err, "could not create new secret networking service")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create new secret networking service")
	}

	s.mqw.NetworkingConfigUpdated(user.OrganizationID)

	return c.JSON(http.StatusOK, ns)
}

func (s *Server) deleteNetworkingService(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	dbNetworkingService, err := s.db.GetNetworkingService(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not get networking service")
		return echo.NewHTTPError(http.StatusBadRequest, "Could not get networking service")
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	if dbNetworkingService.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Networking service does not belong to user's organization")
	}

	err = s.secr.DeleteNetworkingService(uid)
	if err != nil {
		s.log.Error(err, "could not delete networking service from secret")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete networking service from secret")
	}

	err = s.db.DeleteNetworkingService(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not delete networking service from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete networking service from database")
	}

	s.mqw.NetworkingConfigUpdated(user.OrganizationID)

	return c.NoContent(http.StatusNoContent)
}

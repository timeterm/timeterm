package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
)

func (s *Server) getEthernetService(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	secretEthernetConfig, err := s.secr.GetEthernetServiceConfig(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read secret ethernet service")
	}

	apiEthernetConfig := EthernetConfigFrom(secretEthernetConfig, uid)
	return c.JSON(http.StatusOK, apiEthernetConfig)
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

	return c.NoContent(http.StatusNoContent)
}

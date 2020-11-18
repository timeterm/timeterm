package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
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

	err = s.secr.DeleteNetworkingService(uid)
	if err != nil {
		s.log.Error(err, "could not delete networking service")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete networking service")
	}

	return c.NoContent(http.StatusNoContent)
}
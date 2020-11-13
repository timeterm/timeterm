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

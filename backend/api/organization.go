package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

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

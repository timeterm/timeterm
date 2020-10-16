package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
)

func (s *Server) getCurrentUser(c echo.Context) error {
	dbUser, ok := authn.UserFromContext(c)
	if !ok {
		s.log.Error(nil, "user not in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	apiUser := UserFrom(dbUser)
	return c.JSON(http.StatusOK, apiUser)
}

func (s *Server) patchUser(c echo.Context) error {
	userID := c.Param("id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	patchData, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	oldDBUser, err := s.db.GetUserByID(c.Request().Context(), uid)
	if err != nil {
		s.log.Error(err, "could not read user from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read user from database")
	}

	oldAPIUser := UserFrom(oldDBUser)

	jsonUser, err := json.Marshal(oldAPIUser)
	if err != nil {
		s.log.Error(err, "could not marshal the old user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not marshal the old user")
	}

	newJsonUser, err := jsonpatch.MergePatch(jsonUser, patchData)
	if err != nil {
		s.log.Error(err, "could not patch user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch user")
	}

	var newAPIUser User
	err = json.Unmarshal(newJsonUser, &newAPIUser)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched user")
	}
}

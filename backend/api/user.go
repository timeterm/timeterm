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

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	if user.ID != uid {
		return echo.NewHTTPError(http.StatusUnauthorized, "ID mismatch")
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

	newJSONUser, err := jsonpatch.MergePatch(jsonUser, patchData)
	if err != nil {
		s.log.Error(err, "could not patch user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch user")
	}

	var newAPIUser User
	err = json.Unmarshal(newJSONUser, &newAPIUser)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched user")
	}

	newAPIUser.ID = oldAPIUser.ID
	newDBUser := UserToDB(newAPIUser)

	err = s.db.ReplaceUser(c.Request().Context(), newDBUser)
	if err != nil {
		s.log.Error(err, "could not update the user in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update the user in the database")
	}

	return c.JSON(http.StatusOK, newAPIUser)
}

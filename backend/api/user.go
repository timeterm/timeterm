package api

import (
	"net/http"

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

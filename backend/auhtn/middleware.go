package authn

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"gitlab.com/timeterm/timeterm/backend/database"
)

const userEchoContextKey = "gitlab.com/timeterm/timeterm/backend/authn/user"

func UserFromContext(c echo.Context) (database.User, bool) {
	user, ok := c.Get(userEchoContextKey).(database.User)
	return user, ok
}

func AddUserToContext(c echo.Context, u database.User) {
	c.Set(userEchoContextKey, u)
}

func Middleware(db *database.Wrapper, log logr.Logger) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-Api-Key",
		AuthScheme: "",
		Validator: func(key string, c echo.Context) (bool, error) {
			token, err := uuid.Parse(key)
			if err != nil {
				return false, echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
			}

			user, err := db.GetUserByToken(c.Request().Context(), token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}

				log.Error(err, "failed to get user by token")
				return false, echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
			}

			AddUserToContext(c, user)

			return true, nil
		},
	})
}

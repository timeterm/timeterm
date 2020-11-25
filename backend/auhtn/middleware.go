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

const (
	deviceEchoContextKey = "gitlab.com/timeterm/timeterm/backend/authn/device"
	organizationEchoContextKey = "gitlab.com/timeterm/timeterm/backend/authn/organization"
	userEchoContextKey = "gitlab.com/timeterm/timeterm/backend/authn/user"
)

func DeviceFromContext(c echo.Context) (database.Device, bool) {
	dev, ok := c.Get(deviceEchoContextKey).(database.Device)
	return dev, ok
}

func AddDeviceToContext(c echo.Context, d database.Device) {
	c.Set(deviceEchoContextKey, d)
}

func OrganizationFromContext(c echo.Context) (database.Organization, bool) {
	org, ok := c.Get(organizationEchoContextKey).(database.Organization)
	return org, ok
}

func AddOrganizationToContext(c echo.Context, o database.Organization) {
	c.Set(organizationEchoContextKey, o)
}

func UserFromContext(c echo.Context) (database.User, bool) {
	user, ok := c.Get(userEchoContextKey).(database.User)
	return user, ok
}

func AddUserToContext(c echo.Context, u database.User) {
	c.Set(userEchoContextKey, u)
}

func UserLoginMiddleware(db *database.Wrapper, log logr.Logger) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-Api-Key",
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

func DeviceRegistrationLoginMiddleware(db *database.Wrapper, log logr.Logger) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			token, err := uuid.Parse(key)
			if err != nil {
				return false, echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
			}

			organization, err := db.GetOrganizationByDeviceRegistrationToken(c.Request().Context(), token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}

				log.Error(err, "failed to get organization by registration token")
				return false, echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
			}

			AddOrganizationToContext(c, organization)

			return true, nil
		},
	})
}

func DeviceLoginMiddleware(db *database.Wrapper, log logr.Logger) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-Api-Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			token, err := uuid.Parse(key)
			if err != nil {
				return false, echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
			}

			device, err := db.GetDeviceByToken(c.Request().Context(), token)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
				}

				log.Error(err, "failed to get device by token")
				return false, echo.NewHTTPError(http.StatusInternalServerError, "Could not query database")
			}

			AddDeviceToContext(c, device)

			return true, nil
		},
	})
}

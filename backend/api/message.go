package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
)

func (s *Server) getAdminMessage(c echo.Context) error {
	sec, err := strconv.ParseInt(c.Param("sec"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid path parameter 'sec'")
	}

	nanosec, err := strconv.ParseInt(c.Param("nanosec"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid path parameter 'nanosec'")
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	tm := time.Unix(sec, nanosec)
	msg, err := s.db.GetAdminMessage(c.Request().Context(), tm, user.OrganizationID)
	if err != nil {
		s.log.Error(err, "could not get admin message")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read admin message from database")
	}

	decrypted, err := s.msgw.Decrypt(user.OrganizationID, []database.AdminMessage{msg})
	if err != nil {
		s.log.Error(err, "could not decrypt admin message")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read admin message from database")
	}

	return c.JSON(http.StatusOK, adminMessageFrom(decrypted[0]))
}

type getAdminMessagesParams struct {
	FromTimestamp *int64  `query:"fromTimestamp"`
	MaxAmount     *uint64 `query:"maxAmount"`
}

type getAdminMessagesResponse struct {
	Data []AdminMessage `json:"data"`
}

func (s *Server) getAdminMessages(c echo.Context) error {
	var params getAdminMessagesParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	var fromTimestamp *time.Time
	if params.FromTimestamp != nil {
		tm := time.Unix(*params.FromTimestamp, 0)
		fromTimestamp = &tm
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	msgs, err := s.db.GetAdminMessages(c.Request().Context(), database.GetAdminMessagesOpts{
		OrganizationID: user.OrganizationID,
		Limit:          params.MaxAmount,
		FromTimestamp:  fromTimestamp,
	})
	if err != nil {
		s.log.Error(err, "could not get admin messages")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read admin messages from database")
	}

	decrypted, err := s.msgw.Decrypt(user.OrganizationID, msgs)
	if err != nil {
		s.log.Error(err, "could not decrypt admin messages")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read admin messages from database")
	}

	jsonMessages := make([]AdminMessage, 0, len(decrypted))
	for _, msg := range decrypted {
		jsonMessages = append(jsonMessages, adminMessageFrom(msg))
	}

	rsp := getAdminMessagesResponse{Data: jsonMessages}
	return c.JSON(http.StatusOK, rsp)
}

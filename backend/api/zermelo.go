package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/integration/zermelo"
	"gitlab.com/timeterm/timeterm/backend/pkg/jsontypes"
)

type GetZermeloAppointmentsParams struct {
	StartTime jsontypes.UnixTime `query:"startTime"`
	EndTime   jsontypes.UnixTime `query:"endTime"`
}

func (s *Server) getZermeloAppointments(c echo.Context) error {
	s.log.Info("got a getZermeloAppointments request")

	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	student, ok := authn.StudentFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	log := s.log.WithValues(
		"deviceId", dev.ID,
		"studentId", student.ID,
		"organizationId", student.OrganizationID,
	)

	var params GetZermeloAppointmentsParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	if dev.OrganizationID != student.OrganizationID {
		log.Error(nil, "device / user organization ID mismatch")
		return echo.NewHTTPError(http.StatusInternalServerError, "Device / user organization ID mismatch")
	}

	client, err := s.newOrganizationZermeloClient(c.Request().Context(), student.OrganizationID)
	if err != nil {
		log.Error(err, "could not create organization Zermelo client")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not request appointments")
	}

	if !student.ZermeloUser.Valid {
		log.Error(nil, "user has no Zermelo user associated")
		return echo.NewHTTPError(http.StatusUnauthorized, "User has no Zermelo user associated")
	}

	appointments, err := client.GetAppointments(
		c.Request().Context(),
		&zermelo.AppointmentsRequest{
			Start:            params.StartTime.Time(),
			End:              params.EndTime.Time(),
			PossibleStudents: []string{student.ZermeloUser.String},
		},
	)
	if err != nil {
		log.Error(err, "could not get appointments from Zermelo")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not request appointments")
	}

	middle := params.StartTime.Time().Add(params.EndTime.Time().Sub(params.StartTime.Time()) / 2)
	participation, err := client.GetAppointmentParticipations(
		c.Request().Context(),
		&zermelo.AppointmentParticipationsRequest{
			Student: student.ZermeloUser.String,
			Week:    zermelo.YearWeekFromTime(middle),
		},
	)
	if err != nil {
		log.Error(err, "could not get appointment participation from Zermelo")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not request appointments")
	}

	participationByAppointmentInstance := make(map[int][]*zermelo.AppointmentParticipation)
	for _, p := range participation.Response.Data {
		participationByAppointmentInstance[p.AppointmentInstance] = append(
			participationByAppointmentInstance[p.AppointmentInstance], p,
		)
	}

	groupedAppointments := make(map[TimeSpan][]CombinedAppointment)
	for _, appointment := range appointments.Response.Data {
		ts := TimeSpanFromAppointment(appointment)

		for _, ap := range participationByAppointmentInstance[appointment.AppointmentInstance] {
			groupedAppointments[ts] = append(groupedAppointments[ts], CombinedAppointment{
				Appointment:   appointment,
				Participation: ap,
			})
		}
	}

	converted := make([]*ZermeloAppointment, 0)
	for _, group := range groupedAppointments {
		var current *ZermeloAppointment
		var alternatives []*ZermeloAppointment

		for _, apt := range group {
			apiAppointment := apt.ToAPI()

			if apt.Participation.IsStudentEnrolled {
				if current != nil {
					alternatives = append(alternatives, current)
				}
				current = apiAppointment
				continue
			}

			if current == nil ||
				apt.Participation.IsAttendancePlanned ||
				apt.Participation.AttendanceType == zermelo.AttendanceTypeMandatory {
				current = apiAppointment
			} else {
				alternatives = append(alternatives, apiAppointment)
			}
		}

		if current == nil {
			current = new(ZermeloAppointment)
		}
		current.Alternatives = alternatives
		converted = append(converted, current)
	}

	rsp := ZermeloAppointmentsResponse{Data: converted}
	return c.JSON(http.StatusOK, &rsp)
}

type EnrollParams struct {
	UnenrollFromParticipation *int
	EnrollIntoParticipation   *int
}

func EnrollParamsFromRequest(r *http.Request) (EnrollParams, error) {
	var p EnrollParams

	if unenroll := r.URL.Query().Get("unenrollFromParticipation"); unenroll != "" {
		unenrollFrom, err := strconv.Atoi(unenroll)
		if err != nil {
			return p, err
		}
		p.UnenrollFromParticipation = &unenrollFrom
	}

	if enroll := r.URL.Query().Get("enrollIntoParticipation"); enroll != "" {
		enrollInto, err := strconv.Atoi(enroll)
		if err != nil {
			return p, err
		}
		p.EnrollIntoParticipation = &enrollInto
	}

	return p, nil
}

type enrollAction int

const (
	enrollActionNone   enrollAction = 0
	enrollActionEnroll enrollAction = 1 << (iota - 1)
	enrollActionUnenroll
	enrollActionSwitch enrollAction = enrollActionEnroll | enrollActionUnenroll
)

func (s *Server) enrollZermelo(c echo.Context) error {
	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	student, ok := authn.StudentFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	log := s.log.WithValues(
		"deviceId", dev.ID,
		"studentId", student.ID,
		"organizationId", student.OrganizationID,
	)

	if dev.OrganizationID != student.OrganizationID {
		log.Error(nil, "device / user organization ID mismatch")
		return echo.NewHTTPError(http.StatusInternalServerError, "Device / user organization ID mismatch")
	}

	params, err := EnrollParamsFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request data")
	}

	client, err := s.newOrganizationZermeloClient(c.Request().Context(), student.OrganizationID)
	if err != nil {
		log.Error(err, "could not create organization Zermelo client")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not request data from Zermelo")
	}

	if !student.ZermeloUser.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "User has no Zermelo user associated")
	}

	action := enrollActionNone

	canUnenroll := false
	if params.UnenrollFromParticipation != nil {
		upart, err := client.GetAppointmentParticipation(c.Request().Context(), *params.UnenrollFromParticipation)
		if err != nil {
			log.Error(err, "could not get participation to unenroll from")

			var serr zermelo.StatusError
			if errors.As(err, &serr) && serr.Code == http.StatusNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Could not get participation to unenroll from")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get participation to unenroll from")
		}
		if upart.StudentCode != student.ZermeloUser.String {
			log.Error(nil, "Unauthorized to unenroll from participation", "participation", upart)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized to unenroll from participation")
		}
		if !upart.AllowedStudentActions.CanSwitch() {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized to switch participation")
		}
		canUnenroll = upart.AllowedStudentActions == zermelo.AllowedStudentActionsAll

		action |= enrollActionUnenroll
	}

	if params.EnrollIntoParticipation != nil {
		epart, err := client.GetAppointmentParticipation(c.Request().Context(), *params.EnrollIntoParticipation)
		if err != nil {
			log.Error(err, "could not get participation to enroll into")

			var serr zermelo.StatusError
			if errors.As(err, &serr) && serr.Code == http.StatusNotFound {
				return echo.NewHTTPError(http.StatusNotFound, "Could not get participation to enroll into")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not get participation to enroll into")
		}
		if epart.StudentCode != student.ZermeloUser.String {
			log.Error(nil, "Unauthorized to enroll into participation", "participation", epart)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized to enroll into participation")
		}
		if !epart.AllowedStudentActions.CanSwitch() {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized to switch participation")
		}

		action |= enrollActionEnroll
	}

	if action == enrollActionUnenroll && !canUnenroll {
		return echo.NewHTTPError(http.StatusForbidden, "Can not unenroll (can only switch)")
	}

	if params.UnenrollFromParticipation != nil {
		if err = client.ChangeParticipation(c.Request().Context(), &zermelo.ChangeParticipationRequest{
			ParticipationID: *params.UnenrollFromParticipation,
			Enrolled:        false,
		}); err != nil {
			log.Error(err, "could not unenroll")

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not unenroll from participation")
		}
	}

	if params.EnrollIntoParticipation != nil {
		if err = client.ChangeParticipation(c.Request().Context(), &zermelo.ChangeParticipationRequest{
			ParticipationID: *params.EnrollIntoParticipation,
			Enrolled:        true,
		}); err != nil {
			log.Error(err, "could not enroll")

			return echo.NewHTTPError(http.StatusInternalServerError, "Could not enroll into participation")
		}
	}

	return c.NoContent(http.StatusOK)
}

func TimeSpanFromAppointment(a *zermelo.Appointment) TimeSpan {
	return TimeSpan{
		StartUnix: a.Start.Time().Unix(),
		EndUnix:   a.End.Time().Unix(),
	}
}

type CombinedAppointment struct {
	Appointment   *zermelo.Appointment
	Participation *zermelo.AppointmentParticipation
}

func (a *CombinedAppointment) ToAPI() *ZermeloAppointment {
	return &ZermeloAppointment{
		ID:                    a.Appointment.ID,
		ParticipationID:       a.Participation.ID,
		AppointmentInstance:   a.Appointment.AppointmentInstance,
		IsOnline:              a.Participation.IsOnline,
		IsOptional:            a.Participation.IsOptional,
		IsStudentEnrolled:     a.Participation.IsStudentEnrolled,
		IsCanceled:            a.Appointment.IsCanceled,
		StartTimeSlotName:     strconv.Itoa(a.Appointment.StartTimeSlot),
		EndTimeSlotName:       strconv.Itoa(a.Appointment.EndTimeSlot),
		Subjects:              a.Appointment.Subjects,
		Locations:             a.Appointment.Locations,
		Teachers:              a.Appointment.Teachers,
		Groups:                a.Participation.Groups,
		StartTime:             a.Appointment.Start,
		EndTime:               a.Appointment.End,
		Content:               a.Participation.Content,
		AvailableSpace:        a.Participation.AvailableSpace,
		Capacity:              a.Appointment.Capacity,
		AllowedStudentActions: a.Participation.AllowedStudentActions,
	}
}

type TimeSpan struct {
	StartUnix, EndUnix int64
}

func (s *Server) newOrganizationZermeloClient(ctx context.Context, organizationID uuid.UUID) (*zermelo.OrganizationClient, error) {
	token, err := s.secr.GetOrganizationZermeloToken(organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not get Zermelo token for organization: %w", err)
	}

	org, err := s.db.GetOrganization(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve organization: %w", err)
	}

	if org.ZermeloInstitution == "" {
		return nil, errors.New("organization has no Zermelo institution configured")
	}
	if len(token) == 0 {
		return nil, errors.New("organization has no Zermelo token configured")
	}

	return zermelo.NewOrganizationClient(s.log, org.ZermeloInstitution, token)
}

type connectZermeloOrganizationParams struct {
	Token string `json:"token"`
}

func (s *Server) connectZermeloOrganization(c echo.Context) error {
	var params connectZermeloOrganizationParams

	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	err = s.secr.UpsertOrganizationZermeloToken(user.OrganizationID, []byte(params.Token))
	if err != nil {
		s.log.Error(err, "could not upsert organization Zermelo token", "organizationId", user.OrganizationID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not save Zermelo token")
	}

	return c.NoContent(http.StatusOK)
}

package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/integration/zermelo"
)

type GetZermeloAppointmentsParams struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`
}

func (s *Server) getZermeloAppointments(c echo.Context) error {
	dev, ok := authn.DeviceFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	student, ok := authn.StudentFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	log := s.log.WithValues(
		"deviceID", dev.ID,
		"studentID", student.ID,
		"organizationID", student.OrganizationID,
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

	appointments, err := client.GetAppointments(
		c.Request().Context(),
		&zermelo.AppointmentsRequest{
			Start:            params.StartTime,
			End:              params.EndTime,
			PossibleStudents: []string{student.ZermeloUser},
		},
	)
	if err != nil {
		log.Error(err, "could not get appointments from Zermelo")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not request appointments")
	}

	participation, err := client.GetAppointmentParticipations(
		c.Request().Context(),
		&zermelo.AppointmentParticipationsRequest{
			Student: student.ZermeloUser,
			Week:    zermelo.YearWeekFromTime(params.StartTime.Add(params.StartTime.Sub(params.EndTime) / 2)),
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

			if !apt.Appointment.IsOptional {
				converted = append(converted, apiAppointment)
				continue
			}

			if apt.Participation.IsStudentEnrolled {
				if current != nil {
					converted = append(converted, apiAppointment)
					continue
				}
				current = apiAppointment
				continue
			}

			alternatives = append(alternatives, apiAppointment)
		}

		if current == nil {
			converted = append(converted, alternatives...)
		} else {
			current.Alternatives = alternatives
			converted = append(converted, current)
		}
	}

	rsp := ZermeloAppointmentsResponse{Data: converted}
	return c.JSON(http.StatusOK, &rsp)
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
		ID:                  a.Appointment.ID,
		AppointmentInstance: a.Appointment.AppointmentInstance,
		IsOnline:            a.Participation.IsOnline,
		IsOptional:          a.Participation.IsOptional,
		IsStudentEnrolled:   a.Participation.IsStudentEnrolled,
		IsCanceled:          a.Appointment.IsCanceled,
		StartTimeSlotName:   strconv.Itoa(a.Appointment.StartTimeSlot),
		EndTimeSlotName:     strconv.Itoa(a.Appointment.EndTimeSlot),
		Subjects:            a.Appointment.Subjects,
		Locations:           a.Appointment.Locations,
		Teachers:            a.Appointment.Teachers,
		StartTime:           a.Appointment.Start.Time(),
		EndTime:             a.Appointment.End.Time(),
		Content:             a.Participation.Content,
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

	return zermelo.NewOrganizationClient(org.ZermeloInstitution, token)
}

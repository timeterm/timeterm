package zermelo

import (
	"strconv"
	"time"
)

type Time struct {
	Time time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	secs, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		t.Time = time.Unix(secs, 0)
	}
	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AppointmentsResponse struct {
	Response struct {
		Data      []ScheduleWeek `json:"data"`
		Details   string         `json:"details"`
		EndRow    int            `json:"endRow"`
		EventID   int            `json:"eventId"`
		Message   string         `json:"message"`
		StartRow  int            `json:"startRow"`
		Status    int            `json:"status"`
		TotalRows int            `json:"totalRows"`
	} `json:"response"`
}

type ScheduleWeek struct {
	Week         string        `json:"week"`
	User         string        `json:"user"`
	Appointments []Appointment `json:"appointments"`
}

type AppointmentType string

const (
	AppointmentTypeLesson AppointmentType = "lesson"
	AppointmentTypeChoice AppointmentType = "choice"
)

type Appointment struct {
	ID                  int64           `json:"id"`
	Start               Time     `json:"start"`
	End                 Time     `json:"end"`
	Canceled            bool            `json:"cancelled"`
	AppointmentType     AppointmentType `json:"appointmentType`
	Online              bool            `json:"online"`
	Optional            bool            `json:"optional"`
	AppointmentInstance int64           `json:"appointmentInstance"`
	StartTimeSlotName   string          `json:"startTimeSlotName"`
	EndTimeSlotName     string          `json:"endTimeSlotName"`
	Subjects            []string        `json:"subjects"`
	Groups              []string        `json:"groups"`
	Locations           []string        `json:"locations"`
	Teachers            []string        `json:"teachers"`
	ChangeDescription   string          `json:"changeDescription"`
	SchedulerRemark     string          `json:"schedulerRemark"`
	Content             string          `json:"content"`
}
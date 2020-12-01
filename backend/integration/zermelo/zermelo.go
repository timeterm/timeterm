package zermelo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ams = mustLoadLocation("Europe/Amsterdam")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Errorf("integration/zermelo: mustLoadLocation: could not load location %q: %w", name, err))
	}
	return loc
}

type OrganizationClient struct {
	Client  *http.Client
	Token   []byte
	BaseURL *url.URL
}

func NewOrganizationClient(institution string, token []byte) (*OrganizationClient, error) {
	baseURL, err := url.Parse(fmt.Sprintf("https://%s.zportal.nl/api/v3", institution))
	if err != nil {
		return nil, err
	}

	return &OrganizationClient{
		Client: &http.Client{
			Transport: &SetHeaderRoundTripper{
				Key:   "Authorization",
				Value: fmt.Sprintf("Bearer %s", string(token)),
			},
		},
		Token:   token,
		BaseURL: baseURL,
	}, nil
}

type SetHeaderRoundTripper struct {
	Key, Value string
	Next       http.RoundTripper
}

func (t *SetHeaderRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	next := t.Next
	if next == nil {
		next = http.DefaultTransport
	}
	r.Header.Set(t.Key, t.Value)
	return t.Next.RoundTrip(r)
}

type YearWeek struct {
	Year, Week int
}

func YearWeekFromTime(t time.Time) YearWeek {
	y, w := t.ISOWeek()
	return YearWeek{y, w}
}

func (yw YearWeek) String() string {
	return fmt.Sprintf("%d%d", yw.Year, yw.Week)
}

type AppointmentsResponse struct {
	Response AppointmentsResponseData `json:"response"`
}

type AppointmentsResponseData struct {
	ResponseMetadata
	Data []*Appointment `json:"data"`
}

type ResponseMetadata struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Details   string `json:"details"`
	EventID   int    `json:"eventId"`
	StartRow  int    `json:"startRow"`
	EndRow    int    `json:"endRow"`
	TotalRows int    `json:"totalRows"`
}

type AppointmentType string

const (
	AppointmentTypeLesson string = "lesson"
)

type Appointment struct {
	ID                         int             `json:"id"`
	AppointmentInstance        int             `json:"appointmentInstance"`
	Start                      UnixTime        `json:"start"`
	End                        UnixTime        `json:"end"`
	StartTimeSlot              int             `json:"startTimeSlot"`
	EndTimeSlot                int             `json:"endTimeSlot"`
	BranchOfSchool             int             `json:"branchOfSchool"`
	Type                       AppointmentType `json:"type"`
	Teachers                   []string        `json:"teachers"`
	Groups                     []string        `json:"groups"`
	GroupsInDepartments        []int           `json:"groupsInDepartments"`
	Locations                  []string        `json:"locations"`
	LocationsOfBranch          []int           `json:"locationsOfBranch"`
	IsOptional                 bool            `json:"optional"`
	IsValid                    bool            `json:"valid"`
	IsCanceled                 bool            `json:"cancelled"`
	HasTeacherChanged          bool            `json:"teacherChanged"`
	HasGroupChanged            bool            `json:"groupChanged"`
	HasLocationChanged         bool            `json:"locationChanged"`
	HasTimeChanged             bool            `json:"timeChanged"`
	ChangeDescription          string          `json:"changeDescription"`
	SchedulerRemark            string          `json:"schedulerRemark"`
	ChoosableInDepartmentCodes []string        `json:"choosableInDepartmentCodes"`
	Remark                     string          `json:"remark"`
	Subjects                   []string        `json:"subjects"`
}

func structJSONFields(s interface{}) (fields []string) {
	structTy := reflect.TypeOf(s)
	for i := 0; i < structTy.NumField(); i++ {
		field := structTy.Field(i)
		name := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(tag, ",")
			if len(parts) > 0 {
				if parts[0] == "-" {
					continue
				}
				name = parts[0]
			}
		}
		fields = append(fields, name)
	}
	return
}

func appointmentJSONFields() []string {
	return structJSONFields(Appointment{})
}

func appointmentParticipationJSONFields() []string {
	return structJSONFields(AppointmentParticipation{})
}

type UnixTime time.Time

func (t UnixTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixTime) String() string {
	return strconv.FormatInt(time.Time(t).Unix(), 10)
}

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	if t == nil {
		return errors.New("UnixTime is nil")
	}
	ts, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = UnixTime(time.Unix(ts, 0))
	return nil
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(time.Time(t).Unix(), 10)
	return []byte(str), nil
}

type AllowedStudentActions string

const (
	AllowedStudentActionsNone AllowedStudentActions = "none"
	AllowedStudentActionsAll  AllowedStudentActions = "all"
)

type AppointmentParticipation struct {
	ID                    int    `json:"id"`
	AppointmentInstance   int    `json:"appointmentInstance"`
	StudentInDepartment   int    `json:"studentInDepartment"`
	IsOptional            bool   `json:"optional"`
	IsStudentEnrolled     bool   `json:"StudentEnrolled"`
	Content               string `json:"content"`
	IsOnline              bool   `json:"online"`
	IsAttendancePlanned   bool   `json:"plannedAttendance"`
	Capacity              int    `json:"capacity"`
	AllowedStudentActions string `json:"allowedStudentActions"`
	StudentCode           string `json:"studentCode"`
	AvailableSpace        int    `json:"availableSpace"`
}

type AppointmentParticipationsResponse struct {
	Response AppointmentParticipationsResponseData `json:"response"`
}

type AppointmentParticipationsResponseData struct {
	ResponseMetadata
	Data []*AppointmentParticipation `json:"data"`
}

type AppointmentsRequest struct {
	Start time.Time
	End   time.Time

	PossibleStudents []string
}

func (c *OrganizationClient) GetAppointments(
	ctx context.Context,
	req *AppointmentsRequest,
) (*AppointmentsResponse, error) {
	uri := c.BaseURL.ResolveReference(&url.URL{
		Path: "/appointments",
		RawQuery: url.Values{
			"valid":            {"true"},
			"start":            {req.Start.String()},
			"end":              {req.End.String()},
			"fields":           {strings.Join(appointmentJSONFields(), ",")},
			"possibleStudents": {strings.Join(req.PossibleStudents, ",")},
		}.Encode(),
	})

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	hrsp, err := c.Client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("could not do request to Zermelo: %w", err)
	}
	defer func() { _ = hrsp.Body.Close() }()

	var rsp AppointmentsResponse
	if err = json.NewDecoder(hrsp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("could not decode Zermelo response: %w", err)
	}
	return &rsp, nil
}

type AppointmentParticipationsRequest struct {
	Student string
	Week    YearWeek
}

func (c *OrganizationClient) GetAppointmentParticipations(
	ctx context.Context,
	req *AppointmentParticipationsRequest,
) (*AppointmentParticipationsResponse, error) {
	uri := c.BaseURL.ResolveReference(&url.URL{
		Path: "/appointmentparticipations",
		RawQuery: url.Values{
			"student": {req.Student},
			"week":    {req.Week.String()},
			"fields":  {strings.Join(appointmentParticipationJSONFields(), ",")},
		}.Encode(),
	})

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	hrsp, err := c.Client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("could not do request to Zermelo: %w", err)
	}
	defer func() { _ = hrsp.Body.Close() }()

	var rsp AppointmentParticipationsResponse
	if err = json.NewDecoder(hrsp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("could not decode Zermelo response: %w", err)
	}
	return &rsp, nil
}

type ChangeParticipationRequest struct {
	ParticipationID int  `json:"id"`
	Enrolled        bool `json:"studentEnrolled"`
}

type StatusError struct {
	Status int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("got HTTP status %d", e.Status)
}

func (c *OrganizationClient) ChangeParticipation(ctx context.Context, req *ChangeParticipationRequest) error {
	uri := c.BaseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("/appointmentparticipations/%d", req.ParticipationID),
	})

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request body: %w", err)
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	hrsp, err := c.Client.Do(hreq)
	if err != nil {
		return fmt.Errorf("could not do request to Zermelo: %w", err)
	}
	if hrsp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change participation: %w", StatusError{Status: hrsp.StatusCode})
	}

	return nil
}

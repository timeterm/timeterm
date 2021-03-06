package zermelo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"

	"gitlab.com/timeterm/timeterm/backend/messages"
	"gitlab.com/timeterm/timeterm/backend/pkg/jsontypes"
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
	Client         *http.Client
	Token          []byte
	BaseURL        *url.URL
	organizationID uuid.UUID
	log            logr.Logger
	msgw           *messages.Wrapper
}

func NewOrganizationClient(
	log logr.Logger,
	msgw *messages.Wrapper,
	organizationID uuid.UUID,
	institution string,
	token []byte,
) (*OrganizationClient, error) {
	baseURL, err := url.Parse(fmt.Sprintf("https://%s.zportal.nl/api/v3/", institution))
	if err != nil {
		return nil, err
	}

	return &OrganizationClient{
		Client: &http.Client{
			Transport: &SetHeaderRoundTripper{
				log:   log,
				Key:   "User-Agent",
				Value: "Timeterm-Backend/1.0.0",
				Next: &SetHeaderRoundTripper{
					log:   log,
					Key:   "Authorization",
					Value: fmt.Sprintf("Bearer %s", token),
				},
			},
		},
		Token:          token,
		BaseURL:        baseURL,
		organizationID: organizationID,
		log:            log,
		msgw:           msgw,
	}, nil
}

type SetHeaderRoundTripper struct {
	Key, Value string
	Next       http.RoundTripper
	log        logr.Logger
}

func (t *SetHeaderRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	next := t.Next
	if next == nil {
		next = http.DefaultTransport
	}
	r.Header.Set(t.Key, t.Value)
	return next.RoundTrip(r)
}

type YearWeek struct {
	Year, Week int
}

func YearWeekFromTime(t time.Time) YearWeek {
	y, w := t.ISOWeek()
	return YearWeek{y, w}
}

func (yw YearWeek) String() string {
	return fmt.Sprintf("%d%02d", yw.Year, yw.Week)
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
	ID                         int                `json:"id"`
	AppointmentInstance        int                `json:"appointmentInstance"`
	Start                      jsontypes.UnixTime `json:"start"`
	End                        jsontypes.UnixTime `json:"end"`
	StartTimeSlot              int                `json:"startTimeSlot"`
	EndTimeSlot                int                `json:"endTimeSlot"`
	BranchOfSchool             int                `json:"branchOfSchool"`
	Type                       AppointmentType    `json:"type"`
	Teachers                   []string           `json:"teachers"`
	Groups                     []string           `json:"groups"`
	GroupsInDepartments        []int              `json:"groupsInDepartments"`
	Locations                  []string           `json:"locations"`
	LocationsOfBranch          []int              `json:"locationsOfBranch"`
	IsOptional                 *bool              `json:"optional"`
	IsValid                    *bool              `json:"valid"`
	IsCanceled                 *bool              `json:"cancelled"`
	HasTeacherChanged          *bool              `json:"teacherChanged"`
	HasGroupChanged            *bool              `json:"groupChanged"`
	HasLocationChanged         *bool              `json:"locationChanged"`
	HasTimeChanged             *bool              `json:"timeChanged"`
	ChangeDescription          string             `json:"changeDescription"`
	SchedulerRemark            string             `json:"schedulerRemark"`
	ChoosableInDepartmentCodes []string           `json:"choosableInDepartmentCodes"`
	Remark                     string             `json:"remark"`
	Subjects                   []string           `json:"subjects"`
	Capacity                   *int               `json:"capacity"`
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

type AllowedStudentActions string

func (a AllowedStudentActions) CanSwitch() bool {
	return a == AllowedStudentActionsSwitch || a == AllowedStudentActionsAll
}

const (
	AllowedStudentActionsNone   AllowedStudentActions = "none"
	AllowedStudentActionsSwitch AllowedStudentActions = "switch"
	AllowedStudentActionsAll    AllowedStudentActions = "all"
)

type AttendanceType string

const (
	AttendanceTypeNone      AttendanceType = "none"
	AttendanceTypeMandatory AttendanceType = "mandatory"
)

type AppointmentParticipation struct {
	ID                    int                   `json:"id"`
	AppointmentInstance   int                   `json:"appointmentInstance"`
	StudentInDepartment   *int                  `json:"studentInDepartment"`
	IsOptional            *bool                 `json:"optional"`
	IsStudentEnrolled     *bool                 `json:"studentEnrolled"`
	Content               string                `json:"content"`
	IsOnline              *bool                 `json:"online"`
	IsAttendancePlanned   *bool                 `json:"plannedAttendance"`
	Capacity              *int                  `json:"capacity"`
	AllowedStudentActions AllowedStudentActions `json:"allowedStudentActions"`
	StudentCode           string                `json:"studentCode"`
	AvailableSpace        *int                  `json:"availableSpace"`
	Groups                []string              `json:"groups"`
	AttendanceType        AttendanceType        `json:"attendanceType"`
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
		Path: "appointments",
		RawQuery: url.Values{
			"valid":            {"true"},
			"start":            {strconv.FormatInt(req.Start.Unix(), 10)},
			"end":              {strconv.FormatInt(req.End.Unix(), 10)},
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

	if hrsp.StatusCode != http.StatusOK {
		c.logFailedRequest(hreq, hrsp, "Could not retrieve appointments for users %v", req.PossibleStudents)
		return nil, fmt.Errorf("got a response with status code %d (%s)", hrsp.StatusCode, hrsp.Status)
	}

	var rsp AppointmentsResponse
	if err = json.NewDecoder(hrsp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("could not decode Zermelo response: %w", err)
	}
	return &rsp, nil
}

func (c *OrganizationClient) logFailedRequest(hreq *http.Request, hrsp *http.Response, msg string, a ...interface{}) {
	var dumpedReqStr string
	if dumpedReq, err := httputil.DumpRequestOut(hreq, true); err != nil {
		dumpedReqStr = "*** request dump failed: " + err.Error() + " ***"
	} else {
		dumpedReqStr = prefixLines(string(dumpedReq), "\t")
	}

	var dumpedRspStr string
	if dumpedRsp, err := httputil.DumpResponse(hrsp, true); err != nil {
		dumpedRspStr = "*** response dump failed: " + err.Error() + " ***"
	} else {
		dumpedRspStr = prefixLines(string(dumpedRsp), "\t")
	}

	const messageFormat = `request to Zermelo failed

=== request ===
%s

=== response ===
%s
`

	c.msgw.
		Start(c.organizationID).Error().
		Summaryf(msg, a...).
		Messagef(messageFormat, dumpedReqStr, dumpedRspStr).
		Log()
}

func prefixLines(of, with string) string {
	var b strings.Builder
	lines := strings.Split(of, "\n")
	for _, line := range lines {
		b.WriteString(with)
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
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
		Path: "appointmentparticipations",
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

	if hrsp.StatusCode != http.StatusOK {
		c.logFailedRequest(hreq, hrsp, "Could not retrieve appointment participations for user %s", req.Student)
		return nil, fmt.Errorf("got a response with status code %d (%s)", hrsp.StatusCode, hrsp.Status)
	}

	var rsp AppointmentParticipationsResponse
	if err = json.NewDecoder(hrsp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("could not decode Zermelo response: %w", err)
	}
	return &rsp, nil
}

func (c *OrganizationClient) GetAppointmentParticipation(
	ctx context.Context,
	id int,
) (*AppointmentParticipation, error) {
	uri := c.BaseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("appointmentparticipations/%d", id),
		RawQuery: url.Values{
			"fields": {strings.Join(appointmentParticipationJSONFields(), ",")},
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

	if hrsp.StatusCode != http.StatusOK {
		return nil, StatusError{Code: hrsp.StatusCode}
	}

	var rsp AppointmentParticipationsResponse
	if err = json.NewDecoder(hrsp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("could not decode Zermelo response: %w", err)
	}
	if len(rsp.Response.Data) != 1 {
		return nil, fmt.Errorf("unexpected Zermelo response")
	}
	return rsp.Response.Data[0], nil
}

type ChangeParticipationRequest struct {
	ParticipationID int  `json:"id"`
	Enrolled        bool `json:"studentEnrolled"`
}

type StatusError struct {
	Code int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("got HTTP status %d", e.Code)
}

func (c *OrganizationClient) ChangeParticipation(ctx context.Context, req *ChangeParticipationRequest) error {
	uri := c.BaseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("appointmentparticipations/%d", req.ParticipationID),
	})

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request body: %w", err)
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	hrsp, err := c.Client.Do(hreq)
	if err != nil {
		return fmt.Errorf("could not do request to Zermelo: %w", err)
	}
	if hrsp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to change participation: %w", StatusError{Code: hrsp.StatusCode})
	}

	return nil
}

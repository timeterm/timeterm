package zermelo

import (
	"net/url"
	"fmt"
    "net/http"
    "github.com/go-logr/logr"
)

type Client struct {
    log logr.Logger
}

type StudentClient struct {
	client Client
	institution string
	studentCode string
}

type Appointment struct {
    
}

func (sc *StudentClient) getAppointments(token string, year, week int) ([]Appointment, error) {
	baseURLStr := fmt.Sprintf("https://%s.zportal.nl/api/v3/liveschedule", sc.institution)
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}
	q := baseURL.Query()
	q.Set("student", sc.studentCode)
	q.Set("week", fmt.Sprintf("%4d%02d", year, week))
	baseURL.RawQuery = q.Encode()
    resp, err := http.Get(baseURL.String())
    if err != nil {
        return nil, err
    }

    defer func() {
        err = resp.Body.Close()
        if err != nil {
            sc.client.log.Error(err, "Failed to close response body in getAppointments")
        }
    }()
}

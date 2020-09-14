package zermelo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
)

type Client struct {
	log logr.Logger
}

type StudentClient struct {
	client      Client
	institution string
	studentCode string
}

type Appointment struct {
}

type tokenRoundTripper struct {
	token string
	next  http.RoundTripper
}

func (t tokenRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.token)

	return t.next.RoundTrip(r)
}

func withToken(token string, hc *http.Client) *http.Client {
	c := *hc
	c.Transport = tokenRoundTripper{
		token: token,
		next:  c.Transport,
	}

	return &c
}

func (sc *StudentClient) getAppointments(ctx context.Context, token string, year, week int) ([]Appointment, error) {
	baseURL := fmt.Sprintf("https://%s.zportal.nl/api/v3/liveschedule", sc.institution)
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := reqURL.Query()
	q.Set("student", sc.studentCode)
	q.Set("week", fmt.Sprintf("%4d%02d", year, week))
	reqURL.RawQuery = q.Encode()

	req := http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
	}

	resp, err := withToken(token, http.DefaultClient).Do(req.WithContext(ctx))
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

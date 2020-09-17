package zermelo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

type tokenRoundTripper struct {
	token string
	next  http.RoundTripper
}

func (t tokenRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.token)

	return t.next.RoundTrip(r)
}

func (sc *StudentClient) Authenticate(ctx context.Context, authCode string) (*AuthResponse, error) {
	reqURL := fmt.Sprintf("https://%s.zportal.nl/api/v3/oauth/token", sc.institution)

	postValues := url.Values{
		"code":       []string{authCode},
		"grant_type": []string{"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			sc.client.log.Error(err, "integration/zermelo: failed to close response body in (*StudentClient).Authenticate")
		}
	}()

	var authRsp AuthResponse
	err = json.NewDecoder(rsp.Body).Decode(&authRsp)
	if err != nil {
		return nil, err
	}

	return &authRsp, nil
}

func (sc *StudentClient) withToken(token string) *http.Client {
	c := *http.DefaultClient
	c.Transport = tokenRoundTripper{
		token: token,
		next:  http.DefaultTransport,
	}

	return &c
}

func (sc *StudentClient) getAppointments(ctx context.Context, token string, year, week int) (*AppointmentsResponse, error) {
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

	rsp, err := sc.withToken(token).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			sc.client.log.Error(err, "integration/zermelo: failed to close response body in (*StudentClient).getAppointments")
		}
	}()

	var dest AppointmentsResponse
	err = json.NewDecoder(rsp.Body).Decode(&dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

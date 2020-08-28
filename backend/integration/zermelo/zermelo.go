package zermelo

import (
    "net/http"
    "github.com/go-logr/logr"
)

type Client struct {
    log logr.Logger
}

type StudentClient struct {
    client Client
}

type Appointment struct {
    
}

func (c *Client) getAppointments(token string, year, week int) (*[]Appointment, error) {
    resp, err := http.Get("https://{institution}.zportal.nl/api/v3/liveschedule?student={student}&week={week}")
    if err != nil {
        return nil, err
    }

    defer func() {
        err = resp.Body.Close()
        if err != nil {
            return nil, err
        }
    }()
}

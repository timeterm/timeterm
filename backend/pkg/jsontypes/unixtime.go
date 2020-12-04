package jsontypes

import (
	"errors"
	"strconv"
	"time"
)

type UnixTime time.Time

func (t UnixTime) Time() time.Time {
	return time.Time(t)
}

func (t UnixTime) String() string {
	return strconv.FormatInt(time.Time(t).Unix(), 10)
}

func (t *UnixTime) UnmarshalText(b []byte) error {
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

func (t *UnixTime) UnmarshalParam(param string) error {
	return t.UnmarshalText([]byte(param))
}

func (t UnixTime) MarshalText() ([]byte, error) {
	str := strconv.FormatInt(time.Time(t).Unix(), 10)
	return []byte(str), nil
}

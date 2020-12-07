package database

import (
	"errors"
	"fmt"
)

type dbError struct {
	message    string
	underlying error
}

func (e *dbError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.underlying.Error())
}

func (e *dbError) Unwrap() error {
	return e.underlying
}

func (e *dbError) Is(other error) bool {
	var odbe *dbError
	if errors.As(other, &odbe) {
		return e.message == odbe.message
	}
	return false
}

func (e dbError) withUnderlying(cause error) *dbError {
	e.underlying = cause
	return &e
}

var ErrConflict = &dbError{message: "conflict"}

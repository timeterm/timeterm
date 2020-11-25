package jwtpatch

import (
	"time"

	"github.com/nats-io/jwt/v2"
)

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func IntPtr(i int) *int {
	return &i
}

func ClaimTypePtr(c jwt.ClaimType) *jwt.ClaimType {
	return &c
}

func SubjectPtr(s jwt.Subject) *jwt.Subject {
	return &s
}

func ExportTypePtr(t jwt.ExportType) *jwt.ExportType {
	return &t
}

func ResponseTypePtr(t jwt.ExportType) *jwt.ExportType {
	return &t
}

func DurationPtr(d time.Duration) *time.Duration {
	return &d
}

func UintPtr(u uint) *uint {
	return &u
}

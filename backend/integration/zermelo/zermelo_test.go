package zermelo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppointmentJSONFields(t *testing.T) {
	fields := appointmentJSONFields()
	assert.NotEmpty(t, fields)
}

func TestAppointmentParticipationJSONFields(t *testing.T) {
	fields := appointmentParticipationJSONFields()
	assert.NotEmpty(t, fields)
}

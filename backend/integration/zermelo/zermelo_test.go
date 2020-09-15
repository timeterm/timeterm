package zermelo

import (
	"context"
	// "fmt"
	"os"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStudentClient_getAppointments(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	sc := StudentClient{
		client: Client{
			log: zapr.NewLogger(logger),
		},
		institution: os.Getenv("ZERMELO_INSTITUTION"),
		studentCode: os.Getenv("ZERMELO_STUDENT_CODE"),
	}

	token := os.Getenv("ZERMELO_TOKEN")
	_, err := sc.getAppointments(context.Background(), token, 2020, 38)
	assert.NoError(t, err)
}
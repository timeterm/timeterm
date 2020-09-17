package main

import (
	"context"

	"github.com/go-logr/zapr"
	_ "github.com/lib/pq"
	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	sugar := logger.Sugar()

	db, err := database.New("postgres://postgres:postgres@localhost/timeterm?sslmode=disable", log)
	if err != nil {
		sugar.Fatalf("Could not open database: %v", err)
	}

	server := api.NewServer(db)
	sugar.Fatalf("Error running API server: %v", server.Run(context.Background()))
}

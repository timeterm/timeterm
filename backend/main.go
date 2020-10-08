package main

import (
	"context"
	"os"

	"github.com/go-logr/zapr"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func() { _ = logger.Sync() }()

	log := zapr.NewLogger(logger)
	sugar := logger.Sugar()

	db, err := database.New(os.Getenv("DATABASE_URL"), log,
		database.WithJanitor(true),
	)
	if err != nil {
		sugar.Fatalf("could not open database: %v", err)
	}

	server, err := api.NewServer(db, log)
	if err != nil {
		sugar.Fatalf("could not create API server: %v", err)
	}
	sugar.Fatalf("Error running API server: %v", server.Run(context.Background()))
}

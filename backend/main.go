package main

import (
	"context"

	"gitlab.com/timeterm/timeterm/backend/api"
	"gitlab.com/timeterm/timeterm/backend/database"
	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := database.Open("postgres://postgres:postgres@localhost/timeterm?sslmode=disable")
	if err != nil {
		sugar.Fatalf("Could not open database: %v", err)
	}

	server := api.NewServer(db)
	sugar.Fatalf("Error running API server: %v", server.Run(context.Background()))
}

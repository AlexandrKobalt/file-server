package main

import (
	"log"
	"log/slog"

	"github.com/AlexandrKobalt/trip-track/backend/file-server/config"
	"github.com/AlexandrKobalt/trip-track/backend/file-server/internal/app"
	"github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/lifecycle"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	logger := slog.New(slog.Default().Handler())

	a := app.New(cfg, logger)

	if err = lifecycle.Run(a); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log/slog"
	"os"

	"github.com/piotr-woojcik/currencies-service/internal/server"
)

func main() {
	exchangeClient := server.NewExchangeClient()
	s := server.NewServer(exchangeClient)

	if err := s.Start(); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func init() {
	var logLevel slog.Level
	environment := os.Getenv("ENVIRONMENT")
	switch environment {
	case "PRODUCTION":
		logLevel = slog.LevelWarn
	case "DEBUG":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
}

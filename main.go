package main

import (
	"api/api"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute  code", "error", err)
		os.Exit(1)
	}
	slog.Info("all systems offline")
}

func run() error {
	err := godotenv.Load(".env") // using godotenv to use an file .env to envrironment variables
	if err != nil {
		slog.Error("failed to read env", "error", err)
	}

	apiKey := os.Getenv("OMDb_KEY")
	handler := api.NewHandler(apiKey)

	s := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

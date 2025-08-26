package main

import (
	"log/slog"
	"os"
	"strings"

	_ "github.com/animans/REST-API-test-task/docs"
	"github.com/animans/REST-API-test-task/http"
	"github.com/animans/REST-API-test-task/infastructure"
	"github.com/joho/godotenv"
)

// @title Subscriptions REST API
// @version         1.0
// @description     CRUDL по подпискам + summary
// @BasePath        /
// @schemes         http
// @host            localhost:8080

// init ...
func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found", "err", err)
	}
}

// main ...
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(os.ExpandEnv("LOG_LEVEL")),
	}))
	slog.SetDefault(logger)

	repo := infastructure.NewServiceRepoPG()
	err := repo.Open()
	if err != nil {
		slog.Error("repo open failed", "err", err)
		os.Exit(1)
	}
	defer repo.Close()
	api := http.NewHandlers(repo)
	if err := api.Start(); err != nil {
		slog.Error("api start err", "err", err)
		os.Exit(1)
	}
}

func logLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

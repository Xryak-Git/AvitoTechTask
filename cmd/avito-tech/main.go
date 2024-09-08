package main

import (
	"avitoTech/internal/config"
	"avitoTech/internal/storage/sqlite"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.LogLevel)
	log.Info("config loaded", slog.String("log level", cfg.LogLevel))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("cannot init storage", err)
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(level string) *slog.Logger {
	logLevel, err := parseLevel(level)
	if err != nil {
		logLevel = slog.LevelDebug
	}
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
	)

	return log
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

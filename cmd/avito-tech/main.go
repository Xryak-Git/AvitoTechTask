package main

import (
	"avitoTech/internal/app"
	"avitoTech/internal/config"
	"avitoTech/internal/handlers"
	"avitoTech/internal/repo"
	"avitoTech/internal/storage/postgres"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := app.SetupLogger(cfg.LogLevel)

	log.Info("config loaded", "log level", cfg.LogLevel)

	storage, err := postgres.New(cfg.StorageURL)
	if err != nil {
		log.Error("cannot init storage", err)
		os.Exit(1)
	}

	repositories := repo.NewRepos(storage)

	m := http.NewServeMux()
	m.Handle("POST /tender/new", handlers.CreateTender(repositories))
	http.ListenAndServe(":7777", m)

	//_ = storage
}

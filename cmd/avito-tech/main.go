package main

import (
	"avitoTech/internal/app"
	"avitoTech/internal/config"
	v1 "avitoTech/internal/controller/http/v1"
	"avitoTech/internal/repo"
	"avitoTech/internal/service"
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
	services := service.NewServices(repositories)

	log.Info("Initializing handlers and routes...")
	r := v1.NewRouter(services)

	http.ListenAndServe(":8080", r)

}

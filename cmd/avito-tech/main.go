package main

import (
	"avitoTech/internal/app"
	"avitoTech/internal/config"
	"avitoTech/internal/repo"
	"avitoTech/internal/router"
	"avitoTech/internal/service"
	"avitoTech/internal/storage/postgres"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := app.SetupLogger(cfg.LogLevel)

	log.Info("config loaded", "log level", cfg.LogLevel)

	storage, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		log.Error("cannot init storage", err)
		os.Exit(1)
	}

	repositories := repo.NewRepos(storage)
	services := service.NewServices(repositories)

	log.Info("Initializing handlers and routes...")
	r := router.NewRouter(services)

	http.ListenAndServe(":8080", r)

}

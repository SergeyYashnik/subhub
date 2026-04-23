package main

import (
	"Test_EM/internal/config"
	"Test_EM/internal/handler"
	"Test_EM/internal/repository/postgres"
	"Test_EM/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadConfig()

	runMigrations(config.Cfg.DbURL)

	repo, err := postgres.New(config.Cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	subService := service.NewSubscriptionService(repo)
	subHandler := handler.NewSubscriptionHandler(subService)

	appHandler := handler.NewHandler(subHandler)

	srv := &http.Server{
		Addr:    ":" + config.Cfg.AppPort,
		Handler: appHandler.InitRoutes(),
	}

	log.Printf("server started on port %s", config.Cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func runMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("migration init error: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration up error: %v", err)
	}
	log.Println("migrations applied successfully")
}

package main

import (
	"log/slog"
	"net/http"
	"os"

	"gonotes/internal/notes"
	"gonotes/internal/notes/config"
	"gonotes/internal/notes/lib/logger"
	"gonotes/internal/notes/storage/sqlite"

	"github.com/go-chi/chi"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.New(cfg.Env)

	router := chi.NewRouter()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		os.Exit(1)
	}

	handler := notes.NewHandler(log, storage)

	router.Post("/notes", handler.Create)
	router.Get("/notes/{id}", handler.Get)
	router.Delete("/notes/{id}", handler.Delete)
	router.Get("/notes", handler.GetAll)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server error", slog.Any("err", err))
	}

	log.Info("server stopped")
}

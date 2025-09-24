package main

import (
	"log/slog"
	"net/http"
	"os"

	"gonotes/internal/notes"
	"gonotes/internal/notes/config"
	"gonotes/internal/notes/storage/sqlite"

	"github.com/go-chi/chi"
)

func main() {
	cfg := config.MustLoadConfig()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
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

	log.Info("starting server")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("server stopped", slog.Any("err", err))
	}
}

package main

import (
	"log/slog"
	"net/http"
	"os"

	"gonotes/internal/notes"

	"github.com/go-chi/chi"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	router := chi.NewRouter()
	storage := notes.NewStorage()
	handler := notes.NewHandler(log, storage)

	router.Post("/notes", handler.Create)
	router.Get("/notes/{id}", handler.Get)
	router.Delete("/notes/{id}", handler.Delete)
	router.Get("/notes", handler.Get)

	log.Info("starting server")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("server stopped", slog.Any("err", err))
	}
}

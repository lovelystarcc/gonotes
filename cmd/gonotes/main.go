package main

import (
	"net/http"

	"gonotes/internal/notes"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	storage := notes.NewStorage()
	handler := notes.NewHandler(storage)

	router.Post("/notes", handler.CreateNote)
	router.Get("/notes/{id}", handler.GetNote)
	router.Delete("/notes/{id}", handler.DeleteNote)
	router.Get("/notes", handler.GetNotes)

	http.ListenAndServe(":8080", router)
}

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (n *Note) Bind(r *http.Request) error {
	if n.Text == "" {
		return fmt.Errorf("text is required")
	}
	return nil
}

func (n *Note) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

var (
	notes = make(map[int]Note)
	id    = 1
)

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewErrResponse(status int, err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		Message:        err.Error(),
	}
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	if err := render.Bind(r, &n); err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, err))
		return
	}

	n.ID = id
	id++
	notes[n.ID] = n

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &n)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, ok := notes[noteID]
	if !ok {
		render.Render(w, r, NewErrResponse(http.StatusNotFound, fmt.Errorf("note with id %d not found", noteID)))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &n)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	_, ok := notes[noteID]
	if !ok {
		render.Render(w, r, NewErrResponse(http.StatusNotFound, fmt.Errorf("note with id %d not found", noteID)))
		return
	}
	delete(notes, noteID)

	render.Status(r, http.StatusNoContent)
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	list := make([]Note, 0, len(notes))
	for _, note := range notes {
		list = append(list, note)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, list)
}

func main() {
	router := chi.NewRouter()
	router.Post("/notes", createNote)
	router.Get("/notes/{id}", getNote)
	router.Delete("/notes/{id}", deleteNote)
	router.Get("/notes", getNotes)

	http.ListenAndServe(":8080", router)
}

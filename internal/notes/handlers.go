package notes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	if err := render.Bind(r, &n); err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, err))
		return
	}

	created := h.storage.Create(n)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &created)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Get(noteID)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, n)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	err = h.storage.Delete(noteID)
	if err != nil {
		render.Render(w, r, NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	list := h.storage.GetAll()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, list)
}

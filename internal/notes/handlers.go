package notes

import (
	"fmt"
	"gonotes/internal/notes/model"
	"gonotes/internal/notes/storage"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	log     *slog.Logger
	storage storage.NoteRepository
}

func NewHandler(log *slog.Logger, storage storage.NoteRepository) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.create"

	log := h.log.With(
		slog.String("op", op),
	)

	var n model.Note
	if err := render.Bind(r, &n); err != nil {
		log.Error("invalid request body", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusBadRequest, err))
		return
	}

	created, err := h.storage.Create(n)
	if err != nil {
		log.Error("failed to create note", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, created)

	log.Info("note created", slog.Int("id", created.ID))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.get"

	log := h.log.With(
		slog.String("op", op),
	)

	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("invalid id format", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Get(noteID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, n)

	log.Info("note retrieved", slog.Int("id", noteID))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.delete"

	log := h.log.With(
		slog.String("op", op),
	)

	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("invalid id format", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Delete(noteID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, n)

	log.Info("note deleted", slog.Int("id", noteID))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getAll"

	log := h.log.With(
		slog.String("op", op),
	)

	list, err := h.storage.GetAll()
	if err != nil {
		log.Error("failed to get all notes", slog.Any("err", err))
		render.Render(w, r, model.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, list)

	log.Info("notes retrieved", slog.Int("count", len(list)))
}

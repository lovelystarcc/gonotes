package notes

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	log     *slog.Logger
	storage *Storage
}

func NewHandler(log *slog.Logger, storage *Storage) *Handler {
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

	var n Note
	if err := render.Bind(r, &n); err != nil {
		log.Error("invalid request body", slog.Any("err", err))
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, err))
		return
	}

	created := h.storage.Create(n)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &created)

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
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Get(noteID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, NewErrResponse(http.StatusNotFound, err))
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
		render.Render(w, r, NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	err = h.storage.Delete(noteID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusNoContent)

	log.Info("note deleted", slog.Int("id", noteID))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getAll"

	log := h.log.With(
		slog.String("op", op),
	)

	list := h.storage.GetAll()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, list)

	log.Info("notes retrieved", slog.Int("count", len(list)))
}

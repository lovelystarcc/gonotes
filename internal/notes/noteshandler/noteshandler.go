package noteshandler

import (
	"fmt"
	"gonotes/internal/api"
	"gonotes/internal/middleware"
	"gonotes/internal/notes/dto"
	"gonotes/internal/notes/entity"
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
	const op = "notes.handler.create"
	log := h.log.With(slog.String("op", op))

	val := r.Context().Value(middleware.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("unathorized")))
		return
	}

	var req dto.NoteRequest
	if err := render.Bind(r, &req); err != nil {
		log.Error("invalid request body", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusBadRequest, err))
		return
	}

	created, err := h.storage.Create(&entity.Note{UserID: userID, Title: req.Title, Content: req.Content})
	if err != nil {
		log.Error("failed to create note", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, dto.NewNoteResponse(created))

	log.Info("note created", slog.Int("id", created.ID), slog.Int("user_id", userID))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "notes.handler.get"
	log := h.log.With(slog.String("op", op))

	val := r.Context().Value(middleware.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("unauthorized")))
		return
	}

	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("invalid id format", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Get(noteID, userID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusNotFound, err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, dto.NewNoteResponse(n))

	log.Info("note retrieved", slog.Int("id", noteID), slog.Int("user_id", userID))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "notes.handler.delete"
	log := h.log.With(slog.String("op", op))

	val := r.Context().Value(middleware.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("unathorized")))
		return
	}

	idParam := chi.URLParam(r, "id")
	noteID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error("invalid id format", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusBadRequest, fmt.Errorf("invalid id format")))
		return
	}

	n, err := h.storage.Delete(noteID, userID)
	if err != nil {
		log.Error("note not found", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusNotFound, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, dto.NewNoteResponse(n))

	log.Info("note deleted", slog.Int("id", noteID))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "notes.handler.getAll"
	log := h.log.With(slog.String("op", op))

	val := r.Context().Value(middleware.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		log.Error("unauthorized")
		render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("unathorized")))
		return
	}

	list, err := h.storage.GetAll(userID)
	if err != nil {
		log.Error("failed to get all notes", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, dto.NewNoteListResponse(list))

	log.Info("notes retrieved", slog.Int("count", len(list)), slog.Int("user_id", userID))
}

package dto

import (
	"fmt"
	"gonotes/internal/notes/entity"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (n *NoteRequest) Bind(r *http.Request) error {
	if n.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

type NoteResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *NoteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewNoteResponse(n *entity.Note) *NoteResponse {
	return &NoteResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
	}
}

func NewNoteListResponse(notes []entity.Note) []render.Renderer {
	list := make([]render.Renderer, len(notes))
	for i, n := range notes {
		list[i] = NewNoteResponse(&n)
	}
	return list
}

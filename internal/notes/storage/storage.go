package storage

import (
	"errors"
	"gonotes/internal/notes/model"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type NoteRepository interface {
	Create(note model.Note) (*model.Note, error)
	Get(id int) (*model.Note, error)
	Delete(id int) (*model.Note, error)
	GetAll() ([]model.Note, error)
}

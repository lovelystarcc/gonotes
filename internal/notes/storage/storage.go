package storage

import (
	"errors"
	"gonotes/internal/notes/entity"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type NoteRepository interface {
	Create(note *entity.Note) (*entity.Note, error)
	Get(id int, userID int) (*entity.Note, error)
	Delete(id int, userID int) (*entity.Note, error)
	GetAll(userID int) ([]entity.Note, error)
}

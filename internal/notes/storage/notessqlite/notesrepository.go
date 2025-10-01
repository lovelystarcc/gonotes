package notessqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"gonotes/internal/notes/entity"
	"gonotes/internal/notes/storage"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(note *entity.Note) (*entity.Note, error) {
	const op = "storage.sqlite.Create"

	time := time.Now()

	res, err := r.db.Exec("INSERT INTO notes(user_id, title, content, created_at) VALUES (?, ?, ?, ?)", note.UserID, note.Title, note.Content, time)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.Note{
		ID:        int(id),
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: time,
	}, nil
}

func (r *NoteRepository) Get(id int, userID int) (*entity.Note, error) {
	const op = "storage.sqlite.Get"

	var n entity.Note
	err := r.db.QueryRow("SELECT id, title, content FROM notes WHERE id = ? AND user_id = ?", id, userID).Scan(&n.ID, &n.Title, &n.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &n, nil
}

func (r *NoteRepository) Delete(id int, userID int) (*entity.Note, error) {
	const op = "storage.sqlite.Delete"

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var n entity.Note
	err = tx.QueryRow("SELECT id, title, content FROM notes WHERE id = ? AND user_id = ?", id, userID).Scan(&n.ID, &n.Title, &n.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec("DELETE FROM notes WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &n, nil
}

func (r *NoteRepository) GetAll(userID int) ([]entity.Note, error) {
	const op = "storage.sqlite.GetAll"

	rows, err := r.db.Query("SELECT id, title, content, created_at FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var n entity.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

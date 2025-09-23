package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"gonotes/internal/notes/model"
	"gonotes/internal/notes/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY,
		text TEXT NOT NULL
	);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Create(note model.Note) (*model.Note, error) {
	const op = "storage.sqlite.Create"

	res, err := s.db.Exec("INSERT INTO notes(text) VALUES (?)", note.Text)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Note{
		ID:   int(id),
		Text: note.Text,
	}, nil
}

func (s *Storage) Get(id int) (*model.Note, error) {
	const op = "storage.sqlite.Get"

	var n model.Note
	err := s.db.QueryRow("SELECT id, text FROM notes WHERE id = ?", id).Scan(&n.ID, &n.Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &n, nil
}

func (s *Storage) Delete(id int) (*model.Note, error) {
	const op = "storage.sqlite.Delete"

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var n model.Note
	err = tx.QueryRow("SELECT id, text FROM notes WHERE id = ?", id).Scan(&n.ID, &n.Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoteNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &n, nil
}

func (s *Storage) GetAll() ([]model.Note, error) {
	const op = "storage.sqlite.GetAll"

	rows, err := s.db.Query("SELECT id, text FROM notes")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.Text); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

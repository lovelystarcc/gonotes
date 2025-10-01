package authsqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"gonotes/internal/auth"
	"gonotes/internal/auth/entity"

	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) (int, error) {
	const op = "auth.sqlite.Create"

	res, err := r.db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", user.Email, user.Password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, auth.ErrUserExists
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

func (r *UserRepository) Check(id int, user_id int) bool {
	const op = "auth.sqlite.Check"

	var exists int
	err := r.db.QueryRow("SELECT 1 FROM users WHERE id = ? LIMIT 1", id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		return false
	}
	return true
}

func (r *UserRepository) Get(email string) (*entity.User, error) {
	const op = "auth.sqlite.Get"

	var user entity.User
	err := r.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *UserRepository) Delete(id int) error {
	const op = "auth.sqlite.Delete"

	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

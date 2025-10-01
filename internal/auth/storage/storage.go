package storage

import "gonotes/internal/auth/entity"

type UserRepository interface {
	Create(*entity.User) (int, error)
	Get(email string) (*entity.User, error)
	Check(int, int) bool
	Delete(int) error
}

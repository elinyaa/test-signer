package repository

import "github.com/elinyaa/test-signer/internal/domain/entity"

type UserRepository interface {
	Exists(id int) bool
	Upsert(user *entity.User) error
	FindByID(id int) (*entity.User, error)
}

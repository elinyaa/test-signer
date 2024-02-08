package inmem

import (
	"errors"

	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/elinyaa/test-signer/internal/domain/repository"
)

type UserRepository struct {
	users map[int]*entity.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]*entity.User),
	}
}

var _ repository.UserRepository = (*UserRepository)(nil)

var userNotFoundError = errors.New("user not found")

func (r *UserRepository) Exists(id int) bool {
	_, ok := r.users[id]
	return ok
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, userNotFoundError
	}
	return user, nil
}

func (r *UserRepository) Upsert(user *entity.User) error {
	r.users[user.ID] = user
	return nil
}

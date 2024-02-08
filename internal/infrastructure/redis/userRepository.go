package redis

import (
	"context"
	"fmt"

	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/elinyaa/test-signer/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

var _ repository.UserRepository = (*UserRepository)(nil)

func idToKey(id int) string {
	return fmt.Sprint(id)
}

func (r *UserRepository) Exists(id int) bool {
	_, err := r.client.Get(context.Background(), idToKey(id)).Result()
	return err == nil
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	var user entity.User
	err := r.client.Get(context.Background(), idToKey(id)).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Upsert(user *entity.User) error {
	return r.client.Set(context.Background(), idToKey(user.ID), user, 0).Err()
}

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/elinyaa/test-signer/internal/domain/repository"
)

type VerifySignatureUsecase interface {
	VerifySignature(ctx context.Context, user entity.User, signature string) (bool, string, time.Time, error)
}

type verifySignature struct {
	userRepository repository.UserRepository
}

func NewVerifySignature(r repository.UserRepository) *verifySignature {
	return &verifySignature{
		userRepository: r,
	}
}

func (v *verifySignature) VerifySignature(ctx context.Context, user entity.User, signature string) (bool, string, time.Time, error) {
	u, err := v.userRepository.FindByID(user.ID)
	if err != nil {
		return false, "", time.Time{}, err
	}
	verified, answers, dt := u.VerifySignature(signature)
	if !verified {
		return false, "", time.Time{}, errors.New("signature not verified")
	}

	return true, answers, dt, nil
}

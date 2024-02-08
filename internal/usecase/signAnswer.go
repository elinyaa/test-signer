package usecase

import (
	"context"
	"fmt"

	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/elinyaa/test-signer/internal/domain/repository"
)

type SignAnswerUsecase interface {
	SignAnswer(ctx context.Context, userId int, name string, questions string, answers string) (string, error)
}

type signAnswerUC struct {
	userRepository repository.UserRepository
}

var _ SignAnswerUsecase = (*signAnswerUC)(nil)

func NewSignAnswer(r repository.UserRepository) *signAnswerUC {
	return &signAnswerUC{
		userRepository: r,
	}
}

func (s *signAnswerUC) SignAnswer(ctx context.Context, userId int, name string, questions string, answers string) (string, error) {
	if !s.userRepository.Exists(userId) {
		fmt.Println("userId", userId, "does not exist")
		user := entity.NewUser(userId, name)
		err := s.userRepository.Upsert(user)
		if err != nil {
			fmt.Println("error", err)
			return "", err
		}
		fmt.Println("userId", userId, "created", user)
	}

	user, err := s.userRepository.FindByID(userId)
	if err != nil {
		return "", err
	}

	sig := user.AddSignedTest(questions, answers)
	err = s.userRepository.Upsert(user)
	if err != nil {
		return "", err
	}

	return sig, nil
}

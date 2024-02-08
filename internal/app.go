package internal

import (
	"context"
	"log"
	"time"

	"github.com/elinyaa/test-signer/internal/domain/entity"
	"github.com/elinyaa/test-signer/internal/usecase"
)

type app struct {
	logger *log.Logger

	signAnswerUsecase      usecase.SignAnswerUsecase
	verifySignatureUsecase usecase.VerifySignatureUsecase
}

type App interface {
	SignAnswer(ctx context.Context, userId int, name string, questions string, answers string) (string, error)
	VerifySignature(ctx context.Context, user entity.User, signature string) (bool, string, time.Time, error)
}

var _ App = (*app)(nil)

func NewApp(
	l *log.Logger,
	signAnswerUsecase usecase.SignAnswerUsecase,
	verifySignatureUsecase usecase.VerifySignatureUsecase,
) *app {
	return &app{
		logger:                 l,
		signAnswerUsecase:      signAnswerUsecase,
		verifySignatureUsecase: verifySignatureUsecase,
	}
}

func (a *app) SignAnswer(ctx context.Context, userId int, name string, questions string, answers string) (string, error) {
	return a.signAnswerUsecase.SignAnswer(ctx, userId, name, questions, answers)
}

func (a *app) VerifySignature(ctx context.Context, user entity.User, signature string) (bool, string, time.Time, error) {
	return a.verifySignatureUsecase.VerifySignature(ctx, user, signature)
}

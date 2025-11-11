package wallet

import (
	"context"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	walletRepository "github.com/nikolajbotalov/itk_academy_test/internal/repositories/wallet"
	"go.uber.org/zap"
)

type UseCase interface {
	GetByID(ctx context.Context, id string) (*domain.Wallet, error)
	ApplyOperation(ctx context.Context, walletID, operationType string, amount int64) error
}

type useCase struct {
	repo   walletRepository.Repository
	logger *zap.Logger
}

func NewUseCase(repo walletRepository.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, logger: logger}
}

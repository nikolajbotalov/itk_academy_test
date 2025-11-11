package wallet

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	"go.uber.org/zap"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*domain.Wallet, error)
	ApplyOperation(ctx context.Context, id, operationType string, amount int64) error
}

type walletRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewWallet(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return &walletRepository{
		db:     db,
		logger: logger,
	}
}

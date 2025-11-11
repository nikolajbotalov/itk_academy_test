package wallet

import (
	"context"
	"errors"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	"go.uber.org/zap"
)

func (uc *useCase) GetByID(ctx context.Context, id string) (*domain.Wallet, error) {
	uc.logger.Info("Processing request", zap.String("id", id))

	if id == "" {
		uc.logger.Warn("ID is empty", zap.String("wallet_id", id))
		return nil, errors.New("id is required")
	}

	wallet, err := uc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("wallet not found")) {
			uc.logger.Warn("wallet not found", zap.String("wallet_id", id))
			return nil, errors.New("wallet not found")
		}
		uc.logger.Error("failed to get wallet", zap.Error(err), zap.String("wallet_id", id))
		return nil, err
	}

	uc.logger.Info("Successfully processed wallet by id", zap.String("wallet_id", id))
	return wallet, nil
}

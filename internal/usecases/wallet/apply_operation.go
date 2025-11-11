package wallet

import (
	"context"
)

func (uc useCase) ApplyOperation(ctx context.Context, walletID, operationType string, amount int64) error {
	return uc.repo.ApplyOperation(ctx, walletID, operationType, amount)
}

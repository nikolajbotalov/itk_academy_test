package wallet

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

const (
	WITHDRAW = "WITHDRAW"
	DEPOSIT  = "DEPOSIT"
)

func (wr *walletRepository) ApplyOperation(ctx context.Context, walletID, operationType string, amount int64) error {
	tx, err := wr.db.Begin(ctx)
	if err != nil {
		wr.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	var affectedRow int64

	switch operationType {
	case DEPOSIT:
		result, err := tx.Exec(ctx, `
			UPDATE wallets
			SET balance = balance + $1, updated_at = NOW()
			WHERE id = $2
			`, amount, walletID)
		if err != nil {
			wr.logger.Error("Failed to deposit", zap.Error(err))
			return err
		}
		affectedRow = result.RowsAffected()
	case WITHDRAW:
		result, err := tx.Exec(ctx, `
			UPDATE wallets
			SET balance = balance - $1, updated_at = NOW()
			WHERE id = $2 AND balance >= $1
			`, amount, walletID)
		if err != nil {
			wr.logger.Error("Failed to withdraw", zap.Error(err))
			return err
		}
		affectedRow = result.RowsAffected()

		if affectedRow == 0 {
			wr.logger.Error("not enough balance to withdraw",
				zap.String("wallet_id", walletID), zap.Int64("amount", amount))
			return errors.New("not enough balance to withdraw")
		}
	default:
		return errors.New("invalid operation type")
	}

	if err := tx.Commit(ctx); err != nil {
		wr.logger.Error("Failed to commit", zap.Error(err))
		return err
	}

	wr.logger.Info("Operation applied",
		zap.String("wallet_id", walletID),
		zap.String("operation", operationType),
		zap.Int64("amount", amount),
	)

	return nil
}

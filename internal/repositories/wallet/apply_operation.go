package wallet

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

const (
	WITHDRAW = "WITHDRAW"
)

func (wr *walletRepository) ApplyOperation(ctx context.Context, walletID, operationType string, amount int64) error {
	tx, err := wr.db.Begin(ctx)
	if err != nil {
		wr.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	var currentBalance int64
	err = tx.QueryRow(ctx, `
        INSERT INTO wallets (id, balance) 
        VALUES ($1, 0) 
        ON CONFLICT (id) DO UPDATE SET balance = wallets.balance 
        RETURNING balance
    `, walletID).Scan(&currentBalance)

	if err != nil {
		wr.logger.Error("Failed to get/create wallet", zap.Error(err))
		return err
	}

	if operationType == WITHDRAW && currentBalance < amount {
		wr.logger.Error("Current balance is too low", zap.String("id", walletID))
		return errors.New("current balance is too low")
	}

	newBalance := currentBalance + amount
	if operationType == WITHDRAW {
		newBalance = currentBalance - amount
	}

	_, err = tx.Exec(ctx, `
        UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2
    `, newBalance, walletID)
	if err != nil {
		wr.logger.Error("Failed to update balance", zap.Error(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		wr.logger.Error("Failed to commit", zap.Error(err))
		return err
	}

	wr.logger.Info("Operation applied",
		zap.String("wallet_id", walletID),
		zap.String("operation", operationType),
		zap.Int64("old_balance", currentBalance),
		zap.Int64("new_balance", newBalance),
	)

	return nil
}

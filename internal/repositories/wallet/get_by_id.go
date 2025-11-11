package wallet

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	"go.uber.org/zap"
)

func (wr *walletRepository) GetById(ctx context.Context, id string) (*domain.Wallet, error) {
	wr.logger.Info("Fetching wallet by id", zap.String("id", id))

	query, args, err := sq.Select("*").From("wallets").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		wr.logger.Error("Failed to build query", zap.Error(err))
		return nil, err
	}

	row := wr.db.QueryRow(ctx, query, args...)

	var wallet domain.Wallet
	if err = row.Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wr.logger.Error("failed to find wallet by id", zap.String("id", id))
			return nil, errors.New("wallet not found")
		}
		wr.logger.Error("failed to scan row", zap.Error(err))
		return nil, err
	}

	wr.logger.Info("Successfully fetched wallet by id", zap.String("id", id))
	return &wallet, nil
}

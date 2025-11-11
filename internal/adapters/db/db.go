package db

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikolajbotalov/itk_academy_test/internal/config"
	"go.uber.org/zap"
	"time"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*DB, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)

	logger.Info("Connecting to PSQL", zap.String("host", cfg.DB.Host), zap.String("port", cfg.DB.Port))

	var pool *pgxpool.Pool
	ctx := context.Background()

	err := retry.Do(
		func() error {
			var err error
			pool, err = pgxpool.New(ctx, connection)
			if err != nil {
				logger.Warn("Failed to create connection pool", zap.Error(err))
				return fmt.Errorf("create pool: %w", err)
			}

			if err = pool.Ping(ctx); err != nil {
				logger.Warn("Failed to ping db", zap.Error(err))
				pool.Close()
				return fmt.Errorf("ping db: %w", err)
			}

			return nil
		},
		retry.Attempts(cfg.RetryAttempts),
		retry.Delay(2*time.Second),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			logger.Info("Retrying connection", zap.Uint("attempt", n+1), zap.Error(err))
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PSQL after retries: %w", err)
	}

	logger.Info("Successfully connected to PSQL")
	return &DB{
		pool:   pool,
		logger: logger,
	}, nil
}

// Pool возвращает пул соединений для использования в репозиториях
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

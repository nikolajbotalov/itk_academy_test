package db

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/nikolajbotalov/itk_academy_test/internal/config"
	"go.uber.org/zap"
	"os"
	"time"
)

func RunMigrations(cfg *config.Config, logger *zap.Logger) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	logger.Debug("Attempting to connect to database", zap.String("dsn", dsn))

	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		logger.Error("Migrations folder does not exist")
		return fmt.Errorf("migrations folder does not exist")
	}
	logger.Debug("Migrations folder found")

	var m *migrate.Migrate
	var err error

	err = retry.Do(
		func() error {
			m, err = migrate.New("file://migrations", dsn)
			if err != nil {
				logger.Warn("DB not ready, retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(cfg.RetryAttempts),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	logger.Info("Migrations applied successfully")
	return nil
}

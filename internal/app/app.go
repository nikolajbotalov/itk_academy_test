package app

import (
	"fmt"
	"github.com/nikolajbotalov/itk_academy_test/internal/adapters/db"
	"github.com/nikolajbotalov/itk_academy_test/internal/config"
	"github.com/nikolajbotalov/itk_academy_test/internal/logger"
	walletRepository "github.com/nikolajbotalov/itk_academy_test/internal/repositories/wallet"
	walletUseCase "github.com/nikolajbotalov/itk_academy_test/internal/usecases/wallet"
	"go.uber.org/zap"
)

type App struct {
	Server *Server
	Logger *zap.Logger
}

func NewApp() (*App, error) {
	zapLogger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("Failed to create zap logger")
		return nil, err
	}

	cfg := config.LoadConfig(zapLogger)

	if err = db.RunMigrations(cfg, zapLogger); err != nil {
		zapLogger.Error("Failed to run migrations", zap.Error(err))
		return nil, err
	}

	dbInstance, err := db.New(cfg, zapLogger)
	if err != nil {
		zapLogger.Error("Failed to initialize db", zap.Error(err))
		return nil, err
	}

	walletRepo := walletRepository.NewWallet(dbInstance.Pool(), zapLogger)
	walletUseCases := walletUseCase.NewUseCase(walletRepo, zapLogger)

	server := NewServer(cfg, walletUseCases, zapLogger)

	return &App{
		Server: server,
	}, nil
}

func (a *App) Close() {
	a.Logger.Info("Application closed")
}

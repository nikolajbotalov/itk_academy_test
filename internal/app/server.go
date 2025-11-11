package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nikolajbotalov/itk_academy_test/internal/config"
	walletHandlers "github.com/nikolajbotalov/itk_academy_test/internal/handlers/wallet"
	"github.com/nikolajbotalov/itk_academy_test/internal/usecases/wallet"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(cfg *config.Config, walletUseCase wallet.UseCase, logger *zap.Logger) *Server {
	router := gin.New()

	address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	router.Use(ginZapLogger(logger))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ping"})
	})

	walletHandlers.SetupWalletRouter(router, walletUseCase, logger)

	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: router,
		},
		logger: logger,
	}
}

func (s *Server) Run() {
	s.logger.Info("Start server", zap.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.logger.Error("Server start error", zap.Error(err))
	}
}

// возвращает middleware для Gin, который использует Zap для логгирования
func ginZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("Request handled",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("size", c.Writer.Size()),
		)
	}
}

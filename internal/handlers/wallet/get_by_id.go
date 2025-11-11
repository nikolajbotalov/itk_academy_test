package wallet

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nikolajbotalov/itk_academy_test/internal/usecases/wallet"
	"go.uber.org/zap"
	"net/http"
)

func GetById(uc wallet.UseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("wallet_uuid")

		fmt.Printf("id:%s", id)

		walletData, err := uc.GetByID(c.Request.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, errors.New("wallet not found")):
				logger.Info("wallet not found", zap.String("wallet_uuid", id))
				c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			case errors.Is(err, errors.New("empty wallet data")):
				logger.Info("Empty wallet data", zap.String("wallet_uuid", id))
				c.JSON(http.StatusBadRequest, gin.H{"error": "empty wallet data"})
			default:
				logger.Error("Failed to process request", zap.String("wallet_uuid", id), zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		logger.Info("Successfully fetched wallet", zap.String("wallet_uuid", id))
		c.JSON(http.StatusOK, walletData)
	}
}

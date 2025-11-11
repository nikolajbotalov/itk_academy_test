package wallet

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	"github.com/nikolajbotalov/itk_academy_test/internal/usecases/wallet"
	"go.uber.org/zap"
	"net/http"
)

func ApplyOperation(uc wallet.UseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, exists := c.Get("apply_operation_req")
		if !exists {
			logger.Error("did not set apply_operation_req")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		opReq := req.(domain.ApplyOperationRequest)

		err := uc.ApplyOperation(c.Request.Context(), opReq.WalletID, opReq.OperationType, opReq.Amount)
		if err != nil {
			switch {
			case errors.Is(err, errors.New("not enough money")):
				logger.Info("not enough money",
					zap.String("wallet_uuid", opReq.WalletID),
					zap.Int64("amount", opReq.Amount))
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "not enough money",
				})
			case errors.Is(err, errors.New("wallet not found")):
				logger.Info("Wallet not found", zap.String("wallet_uuid", opReq.WalletID))
				c.JSON(http.StatusNotFound, gin.H{
					"error": "wallet not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			}
			return
		}

		logger.Info("Successfully applied operation", zap.String("wallet_id", opReq.WalletID))
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolajbotalov/itk_academy_test/internal/domain"
	"go.uber.org/zap"
	"net/http"
)

func ValidateWalletID(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("wallet_uuid")
		if id == "" {
			logger.Warn("invalid wallet uuid", zap.String("wallet_uuid", id))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet uuid"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateApplyOperation(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.ApplyOperationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request: " + err.Error(),
			})
			c.Abort()
			return
		}

		logger.Info("Validated operation",
			zap.String("wallet_id", req.WalletID),
			zap.String("operation_type", req.OperationType),
			zap.Int64("amount", req.Amount),
		)

		c.Set("apply_operation_req", req)
		c.Next()
	}
}

package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolajbotalov/itk_academy_test/internal/middleware"
	"github.com/nikolajbotalov/itk_academy_test/internal/usecases/wallet"
	"go.uber.org/zap"
)

func SetupWalletRouter(g *gin.Engine, uc wallet.UseCase, logger *zap.Logger) {
	walletRoutes := g.Group("api/v1/wallet")
	{
		walletRoutes.GET("/:wallet_uuid", middleware.ValidateWalletID(logger), GetById(uc, logger))
		walletRoutes.POST("/", middleware.ValidateApplyOperation(logger), ApplyOperation(uc, logger))
	}
}

package wallet

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/wallets")

	// PayOS gọi vào endpoint này, không dùng JWT auth.
	group.POST("/webhook", handler.PayOSWebhook)

	group.GET("/deposit/cancel-return", handler.UpdateWalletCancel)

	group.Use(authMiddleware)
	{
		group.POST("/deposit", handler.CreateDeposit)
		group.GET("/transactions", handler.GetMyTransactions)
	}
}

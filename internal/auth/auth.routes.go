package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := api.Group("/auth")
	{
		group.POST("/register", handler.Register)
		group.GET("/verify", handler.VerifyEmail)
		group.POST("/resend", handler.ResendVerification)
		group.POST("/login", handler.Login)
		group.POST("/logout", authMiddleware, handler.Logout)
		group.POST("/refresh-token", handler.RefreshToken)
		group.POST("/send-reset-password", handler.SendResetPassword)
		group.PATCH("/reset-password", handler.ResetPassword)
	}
}

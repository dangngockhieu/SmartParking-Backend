package rfid_card

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/rfid-cards")
	{
		group.GET("", authMiddleware, adminOnly, handler.FindWithFilters)
		group.GET("/user-email", authMiddleware, adminOnly, handler.GetUserIDByEmail)
		group.GET("/my-cards", authMiddleware, handler.GetMyRfidCard)
		group.GET("/statistics", authMiddleware, adminOnly, handler.GetStatistics)
		group.POST("", authMiddleware, adminOnly, handler.Create)
		group.PATCH("/:id", authMiddleware, adminOnly, handler.Update)
	}
}

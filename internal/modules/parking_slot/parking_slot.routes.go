package parking_slot

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/parking-slots")
	{
		group.POST("", authMiddleware, adminOnly, handler.Create)
		group.GET("/:id", authMiddleware, handler.FindByID)
		group.PATCH("/admin/:id", authMiddleware, adminOnly, handler.AdminUpdateStatus)
		group.POST("/sensor", handler.SensorUpdateStatus)
		group.PATCH("/:id/device", authMiddleware, adminOnly, handler.ChangeDevice)
	}
}

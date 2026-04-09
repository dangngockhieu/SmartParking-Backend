package parking_lot

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/parking-lots")
	group.Use()
	{
		group.POST("", authMiddleware, adminOnly, handler.Create)
		group.GET("", handler.FindAll)
		group.GET("/:id", handler.FindByID)
		group.GET("/:id/gates", authMiddleware, adminOnly, handler.GetGatesByLotID)
		group.PATCH("/:id", authMiddleware, adminOnly, handler.Update)
	}
}

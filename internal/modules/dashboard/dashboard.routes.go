package dashboard

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
	authMiddleware gin.HandlerFunc,
	adminOnly gin.HandlerFunc,
) {
	group := router.Group("/dashboard")
	group.Use(authMiddleware, adminOnly)
	{
		group.GET("/parking-flow", handler.GetParkingFlow)
		group.GET("/revenue-by-month", handler.GetRevenueByMonth)
		group.GET("/revenue-by-day", handler.GetRevenueByDay)
	}
}

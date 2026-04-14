package parking_session

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	group := api.Group("/parking-sessions")
	{
		group.GET("", authMiddleware, adminOnly, handler.FindAll)
		group.GET("/my-sessions", authMiddleware, handler.GetByDate)
	}
}

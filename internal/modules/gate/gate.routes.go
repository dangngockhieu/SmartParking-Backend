package gate

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	// /gates
	gates := api.Group("/gates")
	gates.Use(authMiddleware, adminOnly)
	{
		//Tạo cổng mới
		gates.POST("", handler.Create)
		// Cập nhật thông tin cổng
		gates.PATCH("/:id", handler.Update)
	}
}

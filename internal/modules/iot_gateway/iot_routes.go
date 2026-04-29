package iot_gateway

import "github.com/gin-gonic/gin"

// RegisterRoutes đăng ký các route IoT — không cần JWT, không cần rate limit
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	group := api.Group("/iot")
	{
		group.POST("/camera", handler.CameraPlate)
		group.POST("/rfid", handler.RfidScan)
	}
}

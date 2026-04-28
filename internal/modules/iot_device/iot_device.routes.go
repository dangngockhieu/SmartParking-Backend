package iot_device

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler, authMiddleware, adminOnly gin.HandlerFunc) {
	iotDevices := api.Group("/iot-devices")
	iotDevices.Use(authMiddleware)
	{
		iotDevices.POST("", adminOnly, handler.CreateDevice)
		iotDevices.GET("", adminOnly, handler.FindAllDevices)
		iotDevices.PATCH("/:mac_address", adminOnly, handler.UpdateDevice)
	}
}

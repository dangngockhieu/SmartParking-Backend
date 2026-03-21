package response

import "github.com/gin-gonic/gin"

type SuccessBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	if message == "" {
		message = "Thực hiện thành công"
	}

	c.JSON(statusCode, SuccessBody{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	if message == "" {
		message = "Thực hiện thành công"
	}

	c.JSON(statusCode, SuccessBody{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

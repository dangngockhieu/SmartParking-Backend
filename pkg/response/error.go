package response

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func Error(c *gin.Context, statusCode int, errors []string) {
	if len(errors) == 0 {
		errors = []string{"Internal server error"}
	}

	c.JSON(statusCode, ErrorBody{
		Success: false,
		Errors:  errors,
	})
}

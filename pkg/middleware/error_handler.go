package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.StatusCode, []string{appErr.Message})
			return
		}

		response.Error(c, http.StatusInternalServerError, []string{"Internal server error"})
	}
}

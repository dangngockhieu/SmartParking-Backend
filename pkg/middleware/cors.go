package middleware

import (
	"net/http"
	"strings"

	"backend/configs"

	"github.com/gin-gonic/gin"
)

func CORS(cfg *configs.Config) gin.HandlerFunc {
	allowedOrigins := map[string]struct{}{}

	for _, origin := range strings.Split(cfg.CORSAllowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins[origin] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			if _, ok := allowedOrigins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}, ", "))

		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"ngrok-skip-browser-warning",
		}, ", "))

		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join([]string{
			"Content-Length",
			"Set-Cookie",
		}, ", "))

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

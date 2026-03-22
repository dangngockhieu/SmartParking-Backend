package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/auth/token"
	appErrors "backend/internal/common/errors"
	"backend/internal/modules/user"
	"backend/pkg/database"
)

func Auth(tokenService *token.Service, redis *database.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.Error(appErrors.NewUnauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Error(appErrors.NewUnauthorized("Invalid authorization header"))
			c.Abort()
			return
		}

		rawToken := strings.TrimSpace(parts[1])
		if rawToken == "" {
			c.Error(appErrors.NewUnauthorized("Missing bearer token"))
			c.Abort()
			return
		}

		claims, err := tokenService.VerifyAccessToken(rawToken)
		if err != nil {
			c.Error(appErrors.NewUnauthorized("Invalid or expired access token"))
			c.Abort()
			return
		}

		revoked, err := redis.Exists(redis.Key("auth", "revoked_access", "jti", claims.JTI))
		if err != nil {
			c.Error(appErrors.NewInternal("Cannot validate token state"))
			c.Abort()
			return
		}
		if revoked {
			c.Error(appErrors.NewUnauthorized("Access token has been revoked"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("jti", claims.JTI)
		c.Set("exp", claims.Exp)

		c.Next()
	}
}

func RequireRoles(roles ...user.Role) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[string(role)] = struct{}{}
	}

	return func(c *gin.Context) {
		if len(allowed) == 0 {
			c.Next()
			return
		}

		roleValue, exists := c.Get("role")
		if !exists {
			c.Error(appErrors.NewUnauthorized("Unauthorized"))
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Error(appErrors.NewUnauthorized("Unauthorized"))
			c.Abort()
			return
		}

		if _, ok := allowed[role]; !ok {
			c.Error(appErrors.NewForbidden("Forbidden"))
			c.Abort()
			return
		}

		c.Next()
	}
}

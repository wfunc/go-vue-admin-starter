package middleware

import (
	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		current, ok := appauth.GetCurrentUser(c)
		if !ok {
			response.Error(c, errorx.Unauthorized("missing session"))
			c.Abort()
			return
		}
		if !appauth.HasPermission(current, permission) {
			response.Error(c, errorx.Forbidden("permission denied"))
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"strings"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

func Auth(jwt *appauth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			response.Error(c, errorx.Unauthorized("missing bearer token"))
			c.Abort()
			return
		}
		token := strings.TrimSpace(authorization[7:])
		claims, err := jwt.Parse(token)
		if err != nil {
			response.Error(c, errorx.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}
		appauth.SetCurrentUser(c, appauth.CurrentUser{UserID: claims.UserID, Username: claims.Username, Nickname: claims.Nickname, RoleID: claims.RoleID, RoleCode: claims.RoleCode, Permissions: claims.Permissions})
		c.Next()
	}
}

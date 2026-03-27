package middleware

import (
	"log/slog"

	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", "error", recovered, "path", c.Request.URL.Path)
		response.Error(c, nil)
	})
}

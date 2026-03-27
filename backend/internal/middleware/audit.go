package middleware

import (
	"fmt"
	"strings"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	auditentity "github.com/Wei-Shaw/sub2api/internal/domain/audit/entity"
	auditservice "github.com/Wei-Shaw/sub2api/internal/domain/audit/service"
	"github.com/gin-gonic/gin"
)

func Audit(service *auditservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if service == nil || c.Request.Method == "GET" {
			return
		}
		current, _ := appauth.GetCurrentUser(c)
		module := inferModule(c.Request.URL.Path)
		action := strings.ToLower(c.Request.Method)
		if action == "post" {
			action = "create"
		}
		if action == "put" || action == "patch" {
			action = "update"
		}
		if action == "delete" {
			action = "delete"
		}
		_ = service.Record(c.Request.Context(), auditentity.RecordRequest{
			Username:   current.Username,
			Module:     module,
			Action:     action,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			IP:         c.ClientIP(),
			Detail:     fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		})
	}
}

func inferModule(path string) string {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return strings.ReplaceAll(parts[2], "-", "_")
	}
	return "system"
}

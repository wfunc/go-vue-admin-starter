package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	authhandler "github.com/Wei-Shaw/sub2api/internal/domain/auth/handler"
	audithandler "github.com/Wei-Shaw/sub2api/internal/domain/audit/handler"
	auditservice "github.com/Wei-Shaw/sub2api/internal/domain/audit/service"
	menuhandler "github.com/Wei-Shaw/sub2api/internal/domain/menu/handler"
	rolehandler "github.com/Wei-Shaw/sub2api/internal/domain/role/handler"
	systemhandler "github.com/Wei-Shaw/sub2api/internal/domain/system/handler"
	userhandler "github.com/Wei-Shaw/sub2api/internal/domain/user/handler"
	"github.com/Wei-Shaw/sub2api/internal/infrastructure/cache"
	"github.com/Wei-Shaw/sub2api/internal/infrastructure/db"
	"github.com/Wei-Shaw/sub2api/internal/middleware"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RouterParams struct {
	Mode          string
	AllowedOrigins []string
	Logger        *slog.Logger
	SQLDB         *sql.DB
	Redis         *redis.Client
	JWT           *appauth.Manager
	AuditService  *auditservice.Service
	AuthHandler   *authhandler.Handler
	UserHandler   *userhandler.Handler
	RoleHandler   *rolehandler.Handler
	MenuHandler   *menuhandler.Handler
	SystemHandler *systemhandler.Handler
	AuditHandler  *audithandler.Handler
}

func New(params RouterParams) *gin.Engine {
	gin.SetMode(params.Mode)
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Logger(params.Logger), middleware.Recovery(params.Logger), middleware.CORS(params.AllowedOrigins))

	router.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		result := gin.H{"database": "up", "redis": "disabled"}
		if err := db.Ping(ctx, params.SQLDB); err != nil {
			result["database"] = err.Error()
			c.JSON(http.StatusServiceUnavailable, gin.H{"code": "health_unavailable", "message": "database unavailable", "data": result})
			return
		}
		if params.Redis != nil {
			if err := cache.Ping(ctx, params.Redis); err != nil {
				result["redis"] = err.Error()
				c.JSON(http.StatusServiceUnavailable, gin.H{"code": "health_unavailable", "message": "redis unavailable", "data": result})
				return
			}
			result["redis"] = "up"
		}
		response.Success(c, result)
	})

	api := router.Group("/api/v1")
	publicAuth := api.Group("/auth")
	params.AuthHandler.RegisterPublic(publicAuth)
	api.GET("/system-configs/public", params.SystemHandler.Public)

	protected := api.Group("")
	protected.Use(middleware.Auth(params.JWT), middleware.Audit(params.AuditService))
	params.AuthHandler.RegisterProtected(protected.Group("/auth"))

	users := protected.Group("/users")
	users.GET("", middleware.RequirePermission("user:view"), params.UserHandler.List)
	users.POST("", middleware.RequirePermission("user:create"), params.UserHandler.Create)
	users.PUT(":id", middleware.RequirePermission("user:update"), params.UserHandler.Update)
	users.DELETE(":id", middleware.RequirePermission("user:delete"), params.UserHandler.Delete)

	roles := protected.Group("/roles")
	roles.GET("", middleware.RequirePermission("role:view"), params.RoleHandler.List)
	roles.POST("", middleware.RequirePermission("role:create"), params.RoleHandler.Create)
	roles.PUT(":id", middleware.RequirePermission("role:update"), params.RoleHandler.Update)
	roles.DELETE(":id", middleware.RequirePermission("role:delete"), params.RoleHandler.Delete)

	menus := protected.Group("/menus")
	menus.GET("", middleware.RequirePermission("menu:view"), params.MenuHandler.List)
	menus.POST("", middleware.RequirePermission("menu:create"), params.MenuHandler.Create)
	menus.PUT(":id", middleware.RequirePermission("menu:update"), params.MenuHandler.Update)
	menus.DELETE(":id", middleware.RequirePermission("menu:delete"), params.MenuHandler.Delete)

	configs := protected.Group("/system-configs")
	configs.GET("", middleware.RequirePermission("system_config:view"), params.SystemHandler.List)
	configs.GET("/summary", middleware.RequirePermission("dashboard:view"), params.SystemHandler.Summary)
	configs.POST("", middleware.RequirePermission("system_config:create"), params.SystemHandler.Create)
	configs.PUT(":id", middleware.RequirePermission("system_config:update"), params.SystemHandler.Update)
	configs.DELETE(":id", middleware.RequirePermission("system_config:delete"), params.SystemHandler.Delete)

	audits := protected.Group("/audit-logs")
	audits.GET("", middleware.RequirePermission("audit_log:view"), params.AuditHandler.List)

	return router
}

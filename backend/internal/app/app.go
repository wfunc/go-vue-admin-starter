package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/auth"
	authhandler "github.com/Wei-Shaw/sub2api/internal/domain/auth/handler"
	authservice "github.com/Wei-Shaw/sub2api/internal/domain/auth/service"
	audithandler "github.com/Wei-Shaw/sub2api/internal/domain/audit/handler"
	auditrepository "github.com/Wei-Shaw/sub2api/internal/domain/audit/repository"
	auditservice "github.com/Wei-Shaw/sub2api/internal/domain/audit/service"
	menuhandler "github.com/Wei-Shaw/sub2api/internal/domain/menu/handler"
	menurepository "github.com/Wei-Shaw/sub2api/internal/domain/menu/repository"
	menuservice "github.com/Wei-Shaw/sub2api/internal/domain/menu/service"
	rolehandler "github.com/Wei-Shaw/sub2api/internal/domain/role/handler"
	rolerepository "github.com/Wei-Shaw/sub2api/internal/domain/role/repository"
	roleservice "github.com/Wei-Shaw/sub2api/internal/domain/role/service"
	systemhandler "github.com/Wei-Shaw/sub2api/internal/domain/system/handler"
	systemrepository "github.com/Wei-Shaw/sub2api/internal/domain/system/repository"
	systemservice "github.com/Wei-Shaw/sub2api/internal/domain/system/service"
	userhandler "github.com/Wei-Shaw/sub2api/internal/domain/user/handler"
	userrepository "github.com/Wei-Shaw/sub2api/internal/domain/user/repository"
	userservice "github.com/Wei-Shaw/sub2api/internal/domain/user/service"
	"github.com/Wei-Shaw/sub2api/internal/infrastructure/cache"
	"github.com/Wei-Shaw/sub2api/internal/infrastructure/db"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/server"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config       *config.Config
	Logger       *slog.Logger
	Store        *db.Store
	Redis        *redis.Client
	JWT          *auth.Manager
	Server       *http.Server
	AuditSvc     *auditservice.Service
	routerParams server.RouterParams
}

func New(cfg *config.Config) (*App, error) {
	logFormat := "json"
	if cfg.App.Env == "development" {
		logFormat = "text"
	}
	appLogger := logger.New(logFormat, "info")
	store, err := db.New(cfg.Database)
	if err != nil {
		return nil, err
	}
	redisClient := cache.New(cfg.Redis)
	jwtManager := auth.NewManager(cfg.JWT)

	userRepo := userrepository.New(store.Ent)
	roleRepo := rolerepository.New(store.Ent)
	menuRepo := menurepository.New(store.Ent)
	auditRepo := auditrepository.New(store.Ent)
	systemRepo := systemrepository.New(store.Ent)

	userSvc := userservice.New(userRepo)
	roleSvc := roleservice.New(roleRepo)
	menuSvc := menuservice.New(menuRepo)
	auditSvc := auditservice.New(auditRepo)
	systemSvc := systemservice.New(systemRepo, userSvc, roleSvc, menuSvc, auditSvc)
	authSvc := authservice.New(userSvc, roleSvc, menuSvc, jwtManager)

	params := server.RouterParams{
		Mode:           cfg.Server.Mode,
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		Logger:         appLogger,
		SQLDB:          store.SQL,
		Redis:          redisClient,
		JWT:            jwtManager,
		AuditService:   auditSvc,
		AuthHandler:    authhandler.New(authSvc),
		UserHandler:    userhandler.New(userSvc),
		RoleHandler:    rolehandler.New(roleSvc),
		MenuHandler:    menuhandler.New(menuSvc),
		SystemHandler:  systemhandler.New(systemSvc),
		AuditHandler:   audithandler.New(auditSvc),
	}

	return &App{
		Config:       cfg,
		Logger:       appLogger,
		Store:        store,
		Redis:        redisClient,
		JWT:          jwtManager,
		AuditSvc:     auditSvc,
		routerParams: params,
	}, nil
}

func (a *App) Bootstrap(ctx context.Context) error {
	if err := a.Store.Ent.Schema.Create(ctx); err != nil {
		return err
	}
	return seedStarterData(ctx, a.Store.Ent, a.Config)
}

func (a *App) Start() error {
	if a.Server == nil {
		engine := server.New(a.routerParams)
		a.Server = &http.Server{
			Addr:         a.Config.Server.Address(),
			Handler:      engine,
			ReadTimeout:  a.Config.Server.ReadTimeout,
			WriteTimeout: a.Config.Server.WriteTimeout,
			IdleTimeout:  a.Config.Server.IdleTimeout,
		}
	}
	a.Logger.Info("server started", "addr", a.Server.Addr)
	err := a.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.Server == nil {
		return nil
	}
	return a.Server.Shutdown(ctx)
}

func (a *App) Close() {
	if a.Redis != nil { _ = a.Redis.Close() }
	if a.Store != nil { _ = a.Store.Close() }
}

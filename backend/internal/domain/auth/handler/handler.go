package handler

import (
	"net/http"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	authentity "github.com/Wei-Shaw/sub2api/internal/domain/auth/entity"
	authservice "github.com/Wei-Shaw/sub2api/internal/domain/auth/service"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

type Handler struct { service *authservice.Service }

func New(service *authservice.Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) { rg.POST("/login", h.Login) }

func (h *Handler) RegisterProtected(rg *gin.RouterGroup) {
	rg.GET("/me", h.Me)
	rg.PUT("/profile", h.UpdateProfile)
	rg.POST("/change-password", h.ChangePassword)
	rg.GET("/menus", h.Menus)
	rg.POST("/logout", h.Logout)
}

func (h *Handler) Login(c *gin.Context) {
	var input authentity.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid login payload")); return }
	session, err := h.service.Login(c.Request.Context(), input)
	if err != nil { response.Error(c, err); return }
	response.Success(c, session)
}

func (h *Handler) Me(c *gin.Context) {
	current, ok := appauth.GetCurrentUser(c)
	if !ok { response.Error(c, errorx.Unauthorized("missing session")); return }
	profile, err := h.service.Me(c.Request.Context(), current)
	if err != nil { response.Error(c, err); return }
	response.Success(c, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	current, ok := appauth.GetCurrentUser(c)
	if !ok { response.Error(c, errorx.Unauthorized("missing session")); return }
	var input authentity.UpdateProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid profile payload")); return }
	profile, err := h.service.UpdateProfile(c.Request.Context(), current, input)
	if err != nil { response.Error(c, err); return }
	response.Success(c, profile)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	current, ok := appauth.GetCurrentUser(c)
	if !ok { response.Error(c, errorx.Unauthorized("missing session")); return }
	var input authentity.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid password payload")); return }
	if err := h.service.ChangePassword(c.Request.Context(), current, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "password updated"})
}

func (h *Handler) Menus(c *gin.Context) {
	current, ok := appauth.GetCurrentUser(c)
	if !ok { response.Error(c, errorx.Unauthorized("missing session")); return }
	menus, err := h.service.Menus(c.Request.Context(), current)
	if err != nil { response.Error(c, err); return }
	response.Success(c, menus)
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "success", "message": "logged out"})
}

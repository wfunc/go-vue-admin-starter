package handler

import (
	"net/http"
	"strconv"

	menuentity "github.com/Wei-Shaw/sub2api/internal/domain/menu/entity"
	menuservice "github.com/Wei-Shaw/sub2api/internal/domain/menu/service"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

type Handler struct { service *menuservice.Service }

func New(service *menuservice.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.PUT(":id", h.Update)
	rg.DELETE(":id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), c.Query("keyword"))
	if err != nil { response.Error(c, err); return }
	response.Success(c, items)
}

func (h *Handler) Create(c *gin.Context) {
	var input menuentity.CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid menu payload")); return }
	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil { response.Error(c, err); return }
	response.Created(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { response.Error(c, errorx.BadRequest("invalid menu id")); return }
	var input menuentity.UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid menu payload")); return }
	item, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil { response.Error(c, err); return }
	response.Success(c, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { response.Error(c, errorx.BadRequest("invalid menu id")); return }
	if err := h.service.Delete(c.Request.Context(), id); err != nil { response.Error(c, err); return }
	c.JSON(http.StatusOK, gin.H{"code": "success", "message": "deleted"})
}

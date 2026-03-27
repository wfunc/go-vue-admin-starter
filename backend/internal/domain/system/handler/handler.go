package handler

import (
	"net/http"
	"strconv"

	systementity "github.com/Wei-Shaw/sub2api/internal/domain/system/entity"
	systemservice "github.com/Wei-Shaw/sub2api/internal/domain/system/service"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

type Handler struct { service *systemservice.Service }

func New(service *systemservice.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.PUT(":id", h.Update)
	rg.DELETE(":id", h.Delete)
	rg.GET("/public", h.Public)
	rg.GET("/summary", h.Summary)
}

func (h *Handler) List(c *gin.Context) {
	params := pagination.Parse(c)
	items, total, err := h.service.List(c.Request.Context(), params)
	if err != nil { response.Error(c, err); return }
	response.Paginated(c, items, total, params.Page, params.PageSize)
}

func (h *Handler) Public(c *gin.Context) {
	items, err := h.service.Public(c.Request.Context())
	if err != nil { response.Error(c, err); return }
	response.Success(c, items)
}

func (h *Handler) Summary(c *gin.Context) {
	item, err := h.service.Summary(c.Request.Context())
	if err != nil { response.Error(c, err); return }
	response.Success(c, item)
}

func (h *Handler) Create(c *gin.Context) {
	var input systementity.CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid system config payload")); return }
	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil { response.Error(c, err); return }
	response.Created(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { response.Error(c, errorx.BadRequest("invalid config id")); return }
	var input systementity.UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil { response.Error(c, errorx.BadRequest("invalid system config payload")); return }
	item, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil { response.Error(c, err); return }
	response.Success(c, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { response.Error(c, errorx.BadRequest("invalid config id")); return }
	if err := h.service.Delete(c.Request.Context(), id); err != nil { response.Error(c, err); return }
	c.JSON(http.StatusOK, gin.H{"code": "success", "message": "deleted"})
}

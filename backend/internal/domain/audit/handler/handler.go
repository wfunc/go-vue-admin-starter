package handler

import (
	auditservice "github.com/Wei-Shaw/sub2api/internal/domain/audit/service"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

type Handler struct { service *auditservice.Service }

func New(service *auditservice.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Register(rg *gin.RouterGroup) { rg.GET("", h.List) }

func (h *Handler) List(c *gin.Context) {
	params := pagination.Parse(c)
	items, total, err := h.service.List(c.Request.Context(), params)
	if err != nil { response.Error(c, err); return }
	response.Paginated(c, items, total, params.Page, params.PageSize)
}

package handler

import (
	"strconv"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	conversationentity "github.com/Wei-Shaw/sub2api/internal/domain/conversation/entity"
	conversationservice "github.com/Wei-Shaw/sub2api/internal/domain/conversation/service"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/Wei-Shaw/sub2api/internal/web/response"
	"github.com/gin-gonic/gin"
)

type Handler struct{ service *conversationservice.Service }

func New(service *conversationservice.Service) *Handler { return &Handler{service: service} }

func (h *Handler) List(c *gin.Context) {
	params := pagination.Parse(c)
	status := c.Query("status")
	tier := c.Query("tier")
	keyword := c.Query("keyword")
	items, total, err := h.service.List(c.Request.Context(), params, status, tier, keyword)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Paginated(c, items, total, params.Page, params.PageSize)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.BadRequest("invalid conversation id"))
		return
	}
	item, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) Summary(c *gin.Context) {
	item, err := h.service.Summary(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) Reply(c *gin.Context) {
	id, current, ok := currentUser(c)
	if !ok {
		return
	}
	var input conversationentity.ReplyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errorx.BadRequest("invalid reply payload"))
		return
	}
	if err := h.service.Reply(c.Request.Context(), id, current, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "reply sent"})
}

func (h *Handler) Transfer(c *gin.Context) {
	id, current, ok := currentUser(c)
	if !ok {
		return
	}
	var input conversationentity.TransferRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errorx.BadRequest("invalid transfer payload"))
		return
	}
	if err := h.service.Transfer(c.Request.Context(), id, current, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "conversation transferred"})
}

func (h *Handler) Resolve(c *gin.Context) {
	id, current, ok := currentUser(c)
	if !ok {
		return
	}
	var input conversationentity.ResolveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errorx.BadRequest("invalid resolve payload"))
		return
	}
	if err := h.service.Resolve(c.Request.Context(), id, current, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "conversation resolved"})
}

func currentUser(c *gin.Context) (int, appauth.CurrentUser, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.BadRequest("invalid conversation id"))
		return 0, appauth.CurrentUser{}, false
	}
	current, ok := appauth.GetCurrentUser(c)
	if !ok {
		response.Error(c, errorx.Unauthorized("missing current user"))
		return 0, appauth.CurrentUser{}, false
	}
	return id, current, true
}

package response

import (
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Pagination struct {
	Items    any   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Code: "success", Message: "success", Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Code: "success", Message: "created", Data: data})
}

func Paginated(c *gin.Context, items any, total int64, page, pageSize int) {
	Success(c, Pagination{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func Error(c *gin.Context, err error) {
	appErr := errorx.From(err)
	if appErr == nil {
		appErr = errorx.Internal("internal server error")
	}
	c.JSON(appErr.Status, Envelope{Code: appErr.Code, Message: appErr.Message})
}

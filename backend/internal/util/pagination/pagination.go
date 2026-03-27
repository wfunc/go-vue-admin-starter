package pagination

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page     int
	PageSize int
	Keyword  string
}

func Parse(c *gin.Context) Params {
	page := 1
	pageSize := 10
	if err := bindInt(c.Query("page"), &page); err != nil || page < 1 {
		page = 1
	}
	if err := bindInt(c.Query("page_size"), &pageSize); err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return Params{Page: page, PageSize: pageSize, Keyword: c.Query("keyword")}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func bindInt(raw string, target *int) error {
	if raw == "" {
		return nil
	}
	_, err := fmt.Sscanf(raw, "%d", target)
	return err
}

package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type CreateRequest struct {
	Title      string `json:"title" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	MenuType   string `json:"menu_type"`
	Permission string `json:"permission"`
	Sort       int    `json:"sort"`
	Hidden     bool   `json:"hidden"`
	ParentID   *int   `json:"parent_id"`
}

type UpdateRequest = CreateRequest

type Item struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Component  string    `json:"component"`
	Icon       string    `json:"icon"`
	MenuType   string    `json:"menu_type"`
	Permission string    `json:"permission"`
	Sort       int       `json:"sort"`
	Hidden     bool      `json:"hidden"`
	ParentID   *int      `json:"parent_id,omitempty"`
	Children   []Item    `json:"children,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func FromEnt(model *ent.Menu) Item {
	item := Item{
		ID:         model.ID,
		Title:      model.Title,
		Name:       model.Name,
		Path:       model.Path,
		Component:  model.Component,
		Icon:       model.Icon,
		MenuType:   model.MenuType,
		Permission: model.Permission,
		Sort:       model.Sort,
		Hidden:     model.Hidden,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
	if model.ParentID != nil {
		item.ParentID = model.ParentID
	}
	return item
}

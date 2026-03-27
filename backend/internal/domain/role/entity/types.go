package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type CreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Code        string   `json:"code" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	IsSystem    bool     `json:"is_system"`
}

type UpdateRequest = CreateRequest

type Item struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromEnt(model *ent.Role) Item {
	return Item{
		ID:          model.ID,
		Name:        model.Name,
		Code:        model.Code,
		Description: model.Description,
		Permissions: append([]string{}, model.Permissions...),
		IsSystem:    model.IsSystem,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

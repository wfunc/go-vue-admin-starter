package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type CreateRequest struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateRequest = CreateRequest

type Item struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Summary struct {
	UserCount      int64 `json:"user_count"`
	RoleCount      int64 `json:"role_count"`
	MenuCount      int64 `json:"menu_count"`
	ConfigCount    int64 `json:"config_count"`
	AuditLogCount  int64 `json:"audit_log_count"`
}

func FromEnt(model *ent.SystemConfig) Item {
	return Item{
		ID:          model.ID,
		Key:         model.Key,
		Value:       model.Value,
		Category:    model.Category,
		Description: model.Description,
		IsPublic:    model.IsPublic,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type UpdateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type Item struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Nickname  string    `json:"nickname"`
	Status    string    `json:"status"`
	RoleID    int       `json:"role_id"`
	RoleName  string    `json:"role_name"`
	RoleCode  string    `json:"role_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromEnt(model *ent.User) Item {
	item := Item{
		ID:        model.ID,
		Username:  model.Username,
		Nickname:  model.Nickname,
		Status:    model.Status,
		RoleID:    model.RoleID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	if model.Email != nil {
		item.Email = *model.Email
	}
	if model.Edges.Role != nil {
		item.RoleName = model.Edges.Role.Name
		item.RoleCode = model.Edges.Role.Code
	}
	return item
}

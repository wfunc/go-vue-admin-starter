package entity

import (
	"time"

	menuentity "github.com/Wei-Shaw/sub2api/internal/domain/menu/entity"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type Profile struct {
	UserID      int      `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email,omitempty"`
	Status      string   `json:"status"`
	RoleID      int      `json:"role_id"`
	RoleName    string   `json:"role_name"`
	RoleCode    string   `json:"role_code"`
	Permissions []string `json:"permissions"`
}

type Session struct {
	AccessToken string            `json:"access_token"`
	ExpiresAt   time.Time         `json:"expires_at"`
	User        Profile           `json:"user"`
	Menus       []menuentity.Item `json:"menus"`
}

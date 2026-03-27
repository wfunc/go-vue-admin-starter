package entity

import (
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
)

type RecordRequest struct {
	Username   string `json:"username"`
	Module     string `json:"module"`
	Action     string `json:"action"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	IP         string `json:"ip"`
	Detail     string `json:"detail"`
}

type Item struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Module     string    `json:"module"`
	Action     string    `json:"action"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	IP         string    `json:"ip"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromEnt(model *ent.AuditLog) Item {
	return Item{ID: model.ID, Username: model.Username, Module: model.Module, Action: model.Action, Method: model.Method, Path: model.Path, StatusCode: model.StatusCode, IP: model.IP, Detail: model.Detail, CreatedAt: model.CreatedAt}
}

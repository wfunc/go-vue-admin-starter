package service

import (
	"context"
	"strings"

	auditentity "github.com/Wei-Shaw/sub2api/internal/domain/audit/entity"
	auditrepository "github.com/Wei-Shaw/sub2api/internal/domain/audit/repository"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Service struct { repo *auditrepository.Repository }

func New(repo *auditrepository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Record(ctx context.Context, input auditentity.RecordRequest) error {
	if strings.TrimSpace(input.Username) == "" { input.Username = "anonymous" }
	if strings.TrimSpace(input.Module) == "" { input.Module = "system" }
	if strings.TrimSpace(input.Action) == "" { input.Action = strings.ToLower(input.Method) }
	return s.repo.Record(ctx, input)
}

func (s *Service) List(ctx context.Context, params pagination.Params) ([]auditentity.Item, int64, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil { return nil, 0, err }
	result := make([]auditentity.Item, 0, len(items))
	for _, item := range items { result = append(result, auditentity.FromEnt(item)) }
	return result, total, nil
}

func (s *Service) Count(ctx context.Context) (int64, error) { return s.repo.Count(ctx) }

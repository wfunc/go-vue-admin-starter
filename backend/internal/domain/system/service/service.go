package service

import (
	"context"

	systementity "github.com/Wei-Shaw/sub2api/internal/domain/system/entity"
	systemrepository "github.com/Wei-Shaw/sub2api/internal/domain/system/repository"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Counters interface {
	Count(context.Context) (int64, error)
}

type Service struct {
	repo      *systemrepository.Repository
	users     Counters
	roles     Counters
	menus     Counters
	auditLogs Counters
}

func New(repo *systemrepository.Repository, users Counters, roles Counters, menus Counters, auditLogs Counters) *Service {
	return &Service{repo: repo, users: users, roles: roles, menus: menus, auditLogs: auditLogs}
}

func (s *Service) List(ctx context.Context, params pagination.Params) ([]systementity.Item, int64, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil { return nil, 0, err }
	result := make([]systementity.Item, 0, len(items))
	for _, item := range items { result = append(result, systementity.FromEnt(item)) }
	return result, total, nil
}

func (s *Service) Public(ctx context.Context) (map[string]string, error) {
	items, err := s.repo.Public(ctx)
	if err != nil { return nil, err }
	result := make(map[string]string, len(items))
	for _, item := range items { result[item.Key] = item.Value }
	return result, nil
}

func (s *Service) Create(ctx context.Context, input systementity.CreateRequest) (systementity.Item, error) {
	item, err := s.repo.Create(ctx, input)
	if err != nil { return systementity.Item{}, err }
	return systementity.FromEnt(item), nil
}

func (s *Service) Update(ctx context.Context, id int, input systementity.UpdateRequest) (systementity.Item, error) {
	item, err := s.repo.Update(ctx, id, input)
	if err != nil { return systementity.Item{}, err }
	return systementity.FromEnt(item), nil
}

func (s *Service) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }
func (s *Service) Count(ctx context.Context) (int64, error) { return s.repo.Count(ctx) }

func (s *Service) Summary(ctx context.Context) (systementity.Summary, error) {
	users, err := s.users.Count(ctx)
	if err != nil { return systementity.Summary{}, err }
	roles, err := s.roles.Count(ctx)
	if err != nil { return systementity.Summary{}, err }
	menus, err := s.menus.Count(ctx)
	if err != nil { return systementity.Summary{}, err }
	configs, err := s.repo.Count(ctx)
	if err != nil { return systementity.Summary{}, err }
	audits, err := s.auditLogs.Count(ctx)
	if err != nil { return systementity.Summary{}, err }
	return systementity.Summary{UserCount: users, RoleCount: roles, MenuCount: menus, ConfigCount: configs, AuditLogCount: audits}, nil
}

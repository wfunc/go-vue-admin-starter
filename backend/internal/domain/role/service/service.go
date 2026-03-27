package service

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/domain/role/entity"
	rolerepository "github.com/Wei-Shaw/sub2api/internal/domain/role/repository"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
)

type Service struct {
	repo *rolerepository.Repository
}

func New(repo *rolerepository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, params pagination.Params) ([]entity.Item, int64, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]entity.Item, 0, len(items))
	for _, item := range items {
		result = append(result, entity.FromEnt(item))
	}
	return result, total, nil
}

func (s *Service) Get(ctx context.Context, id int) (entity.Item, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return entity.Item{}, errorx.NotFound("role not found")
	}
	return entity.FromEnt(item), nil
}

func (s *Service) Create(ctx context.Context, input entity.CreateRequest) (entity.Item, error) {
	item, err := s.repo.Create(ctx, input)
	if err != nil {
		return entity.Item{}, err
	}
	return entity.FromEnt(item), nil
}

func (s *Service) Update(ctx context.Context, id int, input entity.UpdateRequest) (entity.Item, error) {
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return entity.Item{}, errorx.NotFound("role not found")
	}
	if current.IsSystem {
		input.IsSystem = true
	}
	item, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return entity.Item{}, err
	}
	return entity.FromEnt(item), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return errorx.NotFound("role not found")
	}
	if current.IsSystem {
		return errorx.Forbidden("system role cannot be deleted")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

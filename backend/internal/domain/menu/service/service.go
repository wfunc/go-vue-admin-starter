package service

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	menuentity "github.com/Wei-Shaw/sub2api/internal/domain/menu/entity"
	menurepository "github.com/Wei-Shaw/sub2api/internal/domain/menu/repository"
)

type Service struct { repo *menurepository.Repository }

func New(repo *menurepository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, keyword string) ([]menuentity.Item, error) {
	items, err := s.repo.List(ctx, keyword)
	if err != nil { return nil, err }
	return buildTree(items), nil
}

func (s *Service) VisibleForPermissions(ctx context.Context, permissions []string) ([]menuentity.Item, error) {
	items, err := s.repo.VisibleForPermissions(ctx, permissions)
	if err != nil { return nil, err }
	return buildTree(items), nil
}

func (s *Service) Create(ctx context.Context, input menuentity.CreateRequest) (menuentity.Item, error) {
	item, err := s.repo.Create(ctx, input)
	if err != nil { return menuentity.Item{}, err }
	return menuentity.FromEnt(item), nil
}

func (s *Service) Update(ctx context.Context, id int, input menuentity.UpdateRequest) (menuentity.Item, error) {
	item, err := s.repo.Update(ctx, id, input)
	if err != nil { return menuentity.Item{}, err }
	return menuentity.FromEnt(item), nil
}

func (s *Service) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }
func (s *Service) Count(ctx context.Context) (int64, error) { return s.repo.Count(ctx) }

func buildTree(models []*dbent.Menu) []menuentity.Item {
	lookup := make(map[int]*menuentity.Item, len(models))
	rootIDs := make([]int, 0)
	for _, model := range models {
		item := menuentity.FromEnt(model)
		item.Children = []menuentity.Item{}
		copy := item
		lookup[item.ID] = &copy
	}
	for _, model := range models {
		current := lookup[model.ID]
		if model.ParentID != nil {
			parent, ok := lookup[*model.ParentID]
			if ok {
				parent.Children = append(parent.Children, *current)
				continue
			}
		}
		rootIDs = append(rootIDs, model.ID)
	}
	roots := make([]menuentity.Item, 0, len(rootIDs))
	for _, id := range rootIDs {
		if item, ok := lookup[id]; ok {
			roots = append(roots, *item)
		}
	}
	return roots
}

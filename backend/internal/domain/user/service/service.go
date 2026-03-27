package service

import (
	"context"
	"strings"

	userentity "github.com/Wei-Shaw/sub2api/internal/domain/user/entity"
	userrepository "github.com/Wei-Shaw/sub2api/internal/domain/user/repository"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
	"github.com/Wei-Shaw/sub2api/internal/util/password"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
)

type Service struct { repo *userrepository.Repository }

func New(repo *userrepository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, params pagination.Params) ([]userentity.Item, int64, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil { return nil, 0, err }
	result := make([]userentity.Item, 0, len(items))
	for _, item := range items { result = append(result, userentity.FromEnt(item)) }
	return result, total, nil
}

func (s *Service) Get(ctx context.Context, id int) (userentity.Item, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil { return userentity.Item{}, errorx.NotFound("user not found") }
	return userentity.FromEnt(item), nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (userentity.Item, string, error) {
	item, err := s.repo.GetByUsername(ctx, username)
	if err != nil { return userentity.Item{}, "", errorx.NotFound("user not found") }
	return userentity.FromEnt(item), item.PasswordHash, nil
}

func (s *Service) Create(ctx context.Context, input userentity.CreateRequest) (userentity.Item, error) {
	hash, err := password.Hash(strings.TrimSpace(input.Password))
	if err != nil { return userentity.Item{}, err }
	item, err := s.repo.Create(ctx, input, hash)
	if err != nil { return userentity.Item{}, err }
	item, err = s.repo.Get(ctx, item.ID)
	if err != nil { return userentity.Item{}, err }
	return userentity.FromEnt(item), nil
}

func (s *Service) Update(ctx context.Context, id int, input userentity.UpdateRequest) (userentity.Item, error) {
	var hash *string
	if strings.TrimSpace(input.Password) != "" {
		value, err := password.Hash(strings.TrimSpace(input.Password))
		if err != nil { return userentity.Item{}, err }
		hash = &value
	}
	item, err := s.repo.Update(ctx, id, input, hash)
	if err != nil { return userentity.Item{}, err }
	item, err = s.repo.Get(ctx, item.ID)
	if err != nil { return userentity.Item{}, err }
	return userentity.FromEnt(item), nil
}

func (s *Service) UpdateProfile(ctx context.Context, id int, nickname, email string) (userentity.Item, error) {
	trimmedNickname := strings.TrimSpace(nickname)
	if trimmedNickname == "" {
		return userentity.Item{}, errorx.BadRequest("nickname is required")
	}
	if _, err := s.repo.Get(ctx, id); err != nil {
		return userentity.Item{}, errorx.NotFound("user not found")
	}
	item, err := s.repo.UpdateProfile(ctx, id, trimmedNickname, email)
	if err != nil { return userentity.Item{}, err }
	item, err = s.repo.Get(ctx, item.ID)
	if err != nil { return userentity.Item{}, err }
	return userentity.FromEnt(item), nil
}

func (s *Service) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	item, err := s.repo.Get(ctx, id)
	if err != nil { return errorx.NotFound("user not found") }
	if err := password.Compare(item.PasswordHash, strings.TrimSpace(currentPassword)); err != nil {
		return errorx.BadRequest("current password is incorrect")
	}
	trimmedNewPassword := strings.TrimSpace(newPassword)
	if len(trimmedNewPassword) < 8 {
		return errorx.BadRequest("new password must be at least 8 characters")
	}
	if err := password.Compare(item.PasswordHash, trimmedNewPassword); err == nil {
		return errorx.BadRequest("new password must be different from current password")
	}
	hash, err := password.Hash(trimmedNewPassword)
	if err != nil { return err }
	return s.repo.UpdatePassword(ctx, id, hash)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if _, err := s.repo.Get(ctx, id); err != nil { return errorx.NotFound("user not found") }
	return s.repo.Delete(ctx, id)
}

func (s *Service) Count(ctx context.Context) (int64, error) { return s.repo.Count(ctx) }

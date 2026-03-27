package service

import (
	"context"
	"strings"

	appauth "github.com/Wei-Shaw/sub2api/internal/auth"
	authentity "github.com/Wei-Shaw/sub2api/internal/domain/auth/entity"
	menuentity "github.com/Wei-Shaw/sub2api/internal/domain/menu/entity"
	menuservice "github.com/Wei-Shaw/sub2api/internal/domain/menu/service"
	roleservice "github.com/Wei-Shaw/sub2api/internal/domain/role/service"
	userentity "github.com/Wei-Shaw/sub2api/internal/domain/user/entity"
	userservice "github.com/Wei-Shaw/sub2api/internal/domain/user/service"
	"github.com/Wei-Shaw/sub2api/internal/util/password"
	"github.com/Wei-Shaw/sub2api/internal/web/errorx"
)

type Service struct {
	users *userservice.Service
	roles *roleservice.Service
	menus *menuservice.Service
	jwt   *appauth.Manager
}

func New(users *userservice.Service, roles *roleservice.Service, menus *menuservice.Service, jwt *appauth.Manager) *Service {
	return &Service{users: users, roles: roles, menus: menus, jwt: jwt}
}

func (s *Service) Login(ctx context.Context, input authentity.LoginRequest) (authentity.Session, error) {
	user, passwordHash, err := s.users.GetByUsername(ctx, strings.TrimSpace(input.Username))
	if err != nil { return authentity.Session{}, errorx.Unauthorized("invalid username or password") }
	if err := password.Compare(passwordHash, strings.TrimSpace(input.Password)); err != nil {
		return authentity.Session{}, errorx.Unauthorized("invalid username or password")
	}
	role, err := s.roles.Get(ctx, user.RoleID)
	if err != nil { return authentity.Session{}, err }
	current := appauth.CurrentUser{UserID: user.ID, Username: user.Username, Nickname: user.Nickname, RoleID: role.ID, RoleCode: role.Code, Permissions: role.Permissions}
	token, expiresAt, err := s.jwt.Generate(current)
	if err != nil { return authentity.Session{}, err }
	menus, err := s.menus.VisibleForPermissions(ctx, role.Permissions)
	if err != nil { return authentity.Session{}, err }
	return authentity.Session{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        toProfile(user, role.Name, role.Code, role.Permissions),
		Menus:       menus,
	}, nil
}

func (s *Service) Me(ctx context.Context, current appauth.CurrentUser) (authentity.Profile, error) {
	user, err := s.users.Get(ctx, current.UserID)
	if err != nil { return authentity.Profile{}, err }
	role, err := s.roles.Get(ctx, user.RoleID)
	if err != nil { return authentity.Profile{}, err }
	return toProfile(user, role.Name, role.Code, role.Permissions), nil
}

func (s *Service) UpdateProfile(ctx context.Context, current appauth.CurrentUser, input authentity.UpdateProfileRequest) (authentity.Profile, error) {
	user, err := s.users.UpdateProfile(ctx, current.UserID, input.Nickname, input.Email)
	if err != nil { return authentity.Profile{}, err }
	role, err := s.roles.Get(ctx, user.RoleID)
	if err != nil { return authentity.Profile{}, err }
	return toProfile(user, role.Name, role.Code, role.Permissions), nil
}

func (s *Service) ChangePassword(ctx context.Context, current appauth.CurrentUser, input authentity.ChangePasswordRequest) error {
	if strings.TrimSpace(input.NewPassword) != strings.TrimSpace(input.ConfirmPassword) {
		return errorx.BadRequest("password confirmation does not match")
	}
	return s.users.ChangePassword(ctx, current.UserID, input.CurrentPassword, input.NewPassword)
}

func (s *Service) Menus(ctx context.Context, current appauth.CurrentUser) ([]menuentity.Item, error) {
	return s.menus.VisibleForPermissions(ctx, current.Permissions)
}

func toProfile(user userentity.Item, roleName, roleCode string, permissions []string) authentity.Profile {
	return authentity.Profile{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Status:      user.Status,
		RoleID:      user.RoleID,
		RoleName:    roleName,
		RoleCode:    roleCode,
		Permissions: permissions,
	}
}

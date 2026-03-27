package repository

import (
	"context"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	entmenu "github.com/Wei-Shaw/sub2api/ent/menu"
	menuentity "github.com/Wei-Shaw/sub2api/internal/domain/menu/entity"
)

type Repository struct { client *dbent.Client }

func New(client *dbent.Client) *Repository { return &Repository{client: client} }

func (r *Repository) List(ctx context.Context, keyword string) ([]*dbent.Menu, error) {
	query := r.client.Menu.Query().Order(entmenu.BySort(), entmenu.ByID())
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		query.Where(entmenu.Or(entmenu.TitleContainsFold(keyword), entmenu.NameContainsFold(keyword), entmenu.PathContainsFold(keyword)))
	}
	return query.All(ctx)
}

func (r *Repository) VisibleForPermissions(ctx context.Context, permissions []string) ([]*dbent.Menu, error) {
	query := r.client.Menu.Query().Where(entmenu.HiddenEQ(false)).Order(entmenu.BySort(), entmenu.ByID())
	if containsWildcard(permissions) {
		return query.All(ctx)
	}
	if len(permissions) == 0 {
		query.Where(entmenu.PermissionEQ(""))
		return query.All(ctx)
	}
	query.Where(entmenu.Or(entmenu.PermissionEQ(""), entmenu.PermissionIn(permissions...)))
	return query.All(ctx)
}

func (r *Repository) Create(ctx context.Context, input menuentity.CreateRequest) (*dbent.Menu, error) {
	builder := r.client.Menu.Create().
		SetTitle(strings.TrimSpace(input.Title)).
		SetName(strings.TrimSpace(input.Name)).
		SetPath(strings.TrimSpace(input.Path)).
		SetComponent(strings.TrimSpace(input.Component)).
		SetIcon(strings.TrimSpace(input.Icon)).
		SetMenuType(defaultMenuType(input.MenuType)).
		SetPermission(strings.TrimSpace(input.Permission)).
		SetSort(input.Sort).
		SetHidden(input.Hidden)
	if input.ParentID != nil { builder.SetParentID(*input.ParentID) }
	return builder.Save(ctx)
}

func (r *Repository) Update(ctx context.Context, id int, input menuentity.UpdateRequest) (*dbent.Menu, error) {
	builder := r.client.Menu.UpdateOneID(id).
		SetTitle(strings.TrimSpace(input.Title)).
		SetName(strings.TrimSpace(input.Name)).
		SetPath(strings.TrimSpace(input.Path)).
		SetComponent(strings.TrimSpace(input.Component)).
		SetIcon(strings.TrimSpace(input.Icon)).
		SetMenuType(defaultMenuType(input.MenuType)).
		SetPermission(strings.TrimSpace(input.Permission)).
		SetSort(input.Sort).
		SetHidden(input.Hidden)
	if input.ParentID != nil { builder.SetParentID(*input.ParentID) } else { builder.ClearParentID() }
	return builder.Save(ctx)
}

func (r *Repository) Delete(ctx context.Context, id int) error { return r.client.Menu.DeleteOneID(id).Exec(ctx) }
func (r *Repository) Count(ctx context.Context) (int64, error) { count, err := r.client.Menu.Query().Count(ctx); return int64(count), err }

func containsWildcard(items []string) bool {
	for _, item := range items { if item == "*" { return true } }
	return false
}

func defaultMenuType(kind string) string {
	if strings.TrimSpace(kind) == "" { return "menu" }
	return strings.TrimSpace(kind)
}

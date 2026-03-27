package repository

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	entrole "github.com/Wei-Shaw/sub2api/ent/role"
	roleentity "github.com/Wei-Shaw/sub2api/internal/domain/role/entity"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Repository struct {
	client *dbent.Client
}

func New(client *dbent.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) List(ctx context.Context, params pagination.Params) ([]*dbent.Role, int64, error) {
	countQuery := r.client.Role.Query()
	listQuery := r.client.Role.Query().Order(entrole.ByCreatedAt(entsql.OrderDesc()))
	applyFilters(countQuery, params.Keyword)
	applyFilters(listQuery, params.Keyword)
	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := listQuery.Offset(params.Offset()).Limit(params.PageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return items, int64(total), nil
}

func (r *Repository) Get(ctx context.Context, id int) (*dbent.Role, error) {
	return r.client.Role.Get(ctx, id)
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*dbent.Role, error) {
	return r.client.Role.Query().Where(entrole.CodeEQ(code)).Only(ctx)
}

func (r *Repository) Create(ctx context.Context, input roleentity.CreateRequest) (*dbent.Role, error) {
	return r.client.Role.Create().
		SetName(input.Name).
		SetCode(strings.TrimSpace(input.Code)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetPermissions(input.Permissions).
		SetIsSystem(input.IsSystem).
		Save(ctx)
}

func (r *Repository) Update(ctx context.Context, id int, input roleentity.UpdateRequest) (*dbent.Role, error) {
	return r.client.Role.UpdateOneID(id).
		SetName(input.Name).
		SetCode(strings.TrimSpace(input.Code)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetPermissions(input.Permissions).
		SetIsSystem(input.IsSystem).
		Save(ctx)
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	return r.client.Role.DeleteOneID(id).Exec(ctx)
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.client.Role.Query().Count(ctx)
	return int64(count), err
}

func applyFilters(query *dbent.RoleQuery, keyword string) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return
	}
	query.Where(entrole.Or(
		entrole.NameContainsFold(keyword),
		entrole.CodeContainsFold(keyword),
		entrole.DescriptionContainsFold(keyword),
	))
}

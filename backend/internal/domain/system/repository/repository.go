package repository

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	systemconfig "github.com/Wei-Shaw/sub2api/ent/systemconfig"
	systementity "github.com/Wei-Shaw/sub2api/internal/domain/system/entity"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Repository struct { client *dbent.Client }

func New(client *dbent.Client) *Repository { return &Repository{client: client} }

func (r *Repository) List(ctx context.Context, params pagination.Params) ([]*dbent.SystemConfig, int64, error) {
	countQuery := r.client.SystemConfig.Query()
	listQuery := r.client.SystemConfig.Query().Order(systemconfig.ByUpdatedAt(entsql.OrderDesc()))
	if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
		predicate := systemconfig.Or(
			systemconfig.KeyContainsFold(keyword),
			systemconfig.CategoryContainsFold(keyword),
			systemconfig.DescriptionContainsFold(keyword),
		)
		countQuery.Where(predicate)
		listQuery.Where(predicate)
	}
	total, err := countQuery.Count(ctx)
	if err != nil { return nil, 0, err }
	items, err := listQuery.Offset(params.Offset()).Limit(params.PageSize).All(ctx)
	if err != nil { return nil, 0, err }
	return items, int64(total), nil
}

func (r *Repository) Public(ctx context.Context) ([]*dbent.SystemConfig, error) {
	return r.client.SystemConfig.Query().Where(systemconfig.IsPublicEQ(true)).All(ctx)
}

func (r *Repository) Create(ctx context.Context, input systementity.CreateRequest) (*dbent.SystemConfig, error) {
	return r.client.SystemConfig.Create().
		SetKey(strings.TrimSpace(input.Key)).
		SetValue(input.Value).
		SetCategory(defaultCategory(input.Category)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetIsPublic(input.IsPublic).
		Save(ctx)
}

func (r *Repository) Update(ctx context.Context, id int, input systementity.UpdateRequest) (*dbent.SystemConfig, error) {
	return r.client.SystemConfig.UpdateOneID(id).
		SetKey(strings.TrimSpace(input.Key)).
		SetValue(input.Value).
		SetCategory(defaultCategory(input.Category)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetIsPublic(input.IsPublic).
		Save(ctx)
}

func (r *Repository) Delete(ctx context.Context, id int) error { return r.client.SystemConfig.DeleteOneID(id).Exec(ctx) }
func (r *Repository) Count(ctx context.Context) (int64, error) { count, err := r.client.SystemConfig.Query().Count(ctx); return int64(count), err }

func defaultCategory(category string) string {
	if strings.TrimSpace(category) == "" { return "general" }
	return strings.TrimSpace(category)
}

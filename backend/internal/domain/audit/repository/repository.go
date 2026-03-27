package repository

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	auditlog "github.com/Wei-Shaw/sub2api/ent/auditlog"
	auditentity "github.com/Wei-Shaw/sub2api/internal/domain/audit/entity"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Repository struct { client *dbent.Client }

func New(client *dbent.Client) *Repository { return &Repository{client: client} }

func (r *Repository) Record(ctx context.Context, input auditentity.RecordRequest) error {
	_, err := r.client.AuditLog.Create().
		SetUsername(input.Username).
		SetModule(input.Module).
		SetAction(input.Action).
		SetMethod(input.Method).
		SetPath(input.Path).
		SetStatusCode(input.StatusCode).
		SetIP(input.IP).
		SetDetail(input.Detail).
		Save(ctx)
	return err
}

func (r *Repository) List(ctx context.Context, params pagination.Params) ([]*dbent.AuditLog, int64, error) {
	countQuery := r.client.AuditLog.Query()
	listQuery := r.client.AuditLog.Query().Order(auditlog.ByCreatedAt(entsql.OrderDesc()))
	if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
		predicate := auditlog.Or(auditlog.UsernameContainsFold(keyword), auditlog.PathContainsFold(keyword), auditlog.ActionContainsFold(keyword), auditlog.ModuleContainsFold(keyword))
		countQuery.Where(predicate)
		listQuery.Where(predicate)
	}
	total, err := countQuery.Count(ctx)
	if err != nil { return nil, 0, err }
	items, err := listQuery.Offset(params.Offset()).Limit(params.PageSize).All(ctx)
	if err != nil { return nil, 0, err }
	return items, int64(total), nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) { count, err := r.client.AuditLog.Query().Count(ctx); return int64(count), err }

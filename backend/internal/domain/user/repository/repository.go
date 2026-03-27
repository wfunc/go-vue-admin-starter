package repository

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	entuser "github.com/Wei-Shaw/sub2api/ent/user"
	userentity "github.com/Wei-Shaw/sub2api/internal/domain/user/entity"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Repository struct { client *dbent.Client }

func New(client *dbent.Client) *Repository { return &Repository{client: client} }

func (r *Repository) List(ctx context.Context, params pagination.Params) ([]*dbent.User, int64, error) {
	countQuery := r.client.User.Query()
	listQuery := r.client.User.Query().WithRole().Order(entuser.ByCreatedAt(entsql.OrderDesc()))
	applyFilters(countQuery, params.Keyword)
	applyFilters(listQuery, params.Keyword)
	total, err := countQuery.Count(ctx)
	if err != nil { return nil, 0, err }
	items, err := listQuery.Offset(params.Offset()).Limit(params.PageSize).All(ctx)
	if err != nil { return nil, 0, err }
	return items, int64(total), nil
}

func (r *Repository) Get(ctx context.Context, id int) (*dbent.User, error) {
	return r.client.User.Query().Where(entuser.IDEQ(id)).WithRole().Only(ctx)
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*dbent.User, error) {
	return r.client.User.Query().Where(entuser.UsernameEQ(username)).WithRole().Only(ctx)
}

func (r *Repository) Create(ctx context.Context, input userentity.CreateRequest, passwordHash string) (*dbent.User, error) {
	builder := r.client.User.Create().
		SetUsername(strings.TrimSpace(input.Username)).
		SetPasswordHash(passwordHash).
		SetNickname(strings.TrimSpace(input.Nickname)).
		SetStatus(defaultStatus(input.Status)).
		SetRoleID(input.RoleID)
	if email := strings.TrimSpace(input.Email); email != "" {
		builder.SetEmail(email)
	}
	return builder.Save(ctx)
}

func (r *Repository) Update(ctx context.Context, id int, input userentity.UpdateRequest, passwordHash *string) (*dbent.User, error) {
	builder := r.client.User.UpdateOneID(id).
		SetNickname(strings.TrimSpace(input.Nickname)).
		SetStatus(defaultStatus(input.Status)).
		SetRoleID(input.RoleID)
	if email := strings.TrimSpace(input.Email); email != "" {
		builder.SetEmail(email)
	} else {
		builder.ClearEmail()
	}
	if passwordHash != nil {
		builder.SetPasswordHash(*passwordHash)
	}
	return builder.Save(ctx)
}

func (r *Repository) UpdateProfile(ctx context.Context, id int, nickname, email string) (*dbent.User, error) {
	builder := r.client.User.UpdateOneID(id).SetNickname(strings.TrimSpace(nickname))
	if trimmed := strings.TrimSpace(email); trimmed != "" {
		builder.SetEmail(trimmed)
	} else {
		builder.ClearEmail()
	}
	return builder.Save(ctx)
}

func (r *Repository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	return r.client.User.UpdateOneID(id).SetPasswordHash(passwordHash).Exec(ctx)
}

func (r *Repository) Delete(ctx context.Context, id int) error { return r.client.User.DeleteOneID(id).Exec(ctx) }

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.client.User.Query().Count(ctx)
	return int64(count), err
}

func applyFilters(query *dbent.UserQuery, keyword string) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" { return }
	query.Where(entuser.Or(
		entuser.UsernameContainsFold(keyword),
		entuser.NicknameContainsFold(keyword),
		entuser.EmailContainsFold(keyword),
	))
}

func defaultStatus(status string) string {
	if strings.TrimSpace(status) == "" { return "active" }
	return strings.TrimSpace(status)
}

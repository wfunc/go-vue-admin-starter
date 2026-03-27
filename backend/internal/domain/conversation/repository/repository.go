package repository

import (
	"context"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	entconversation "github.com/Wei-Shaw/sub2api/ent/conversation"
	entconversationmessage "github.com/Wei-Shaw/sub2api/ent/conversationmessage"
	"github.com/Wei-Shaw/sub2api/internal/util/pagination"
)

type Repository struct{ client *dbent.Client }

func New(client *dbent.Client) *Repository { return &Repository{client: client} }

func (r *Repository) List(ctx context.Context, params pagination.Params, status, tier, keyword string) ([]*dbent.Conversation, int64, error) {
	countQuery := r.client.Conversation.Query()
	listQuery := r.client.Conversation.Query().Order(entconversation.ByLastActiveAt(entsql.OrderDesc()))
	applyFilters(countQuery, status, tier, keyword)
	applyFilters(listQuery, status, tier, keyword)
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

func (r *Repository) Get(ctx context.Context, id int) (*dbent.Conversation, error) {
	return r.client.Conversation.Query().
		Where(entconversation.IDEQ(id)).
		WithMessages(func(query *dbent.ConversationMessageQuery) {
			query.Order(entconversationmessage.ByCreatedAt())
		}).
		Only(ctx)
}

func (r *Repository) All(ctx context.Context) ([]*dbent.Conversation, error) {
	return r.client.Conversation.Query().
		Order(entconversation.ByLastActiveAt(entsql.OrderDesc())).
		WithMessages(func(query *dbent.ConversationMessageQuery) {
			query.Order(entconversationmessage.ByCreatedAt(entsql.OrderDesc()))
		}).
		All(ctx)
}

func (r *Repository) Reply(ctx context.Context, id int, actor, actorType, messageType, content string) error {
	now := time.Now()
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	if _, err = tx.ConversationMessage.Create().
		SetConversationID(id).
		SetActor(actor).
		SetActorType(actorType).
		SetMessageType(messageType).
		SetContent(content).
		SetCreatedAt(now).
		Save(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Conversation.UpdateOneID(id).
		SetPreview(content).
		SetLastActiveAt(now).
		SetStatus("processing").
		SetUnread(0).
		ClearClosedAt().
		Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) Transfer(ctx context.Context, id int, actor, assignee, queue, content string) error {
	now := time.Now()
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	if _, err = tx.ConversationMessage.Create().
		SetConversationID(id).
		SetActor(actor).
		SetActorType("agent").
		SetMessageType("event").
		SetContent(content).
		SetCreatedAt(now).
		Save(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	updater := tx.Conversation.UpdateOneID(id).
		SetLastActiveAt(now).
		SetStatus("processing")
	if strings.TrimSpace(assignee) != "" {
		updater.SetAssignee(strings.TrimSpace(assignee))
	}
	if strings.TrimSpace(queue) != "" {
		updater.SetQueue(strings.TrimSpace(queue))
	}
	if err = updater.Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) Resolve(ctx context.Context, id int, actor, preview, content string) error {
	now := time.Now()
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	if _, err = tx.ConversationMessage.Create().
		SetConversationID(id).
		SetActor(actor).
		SetActorType("agent").
		SetMessageType("event").
		SetContent(content).
		SetCreatedAt(now).
		Save(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Conversation.UpdateOneID(id).
		SetStatus("closed").
		SetPreview(preview).
		SetLastActiveAt(now).
		SetClosedAt(now).
		SetUnread(0).
		Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.client.Conversation.Query().Count(ctx)
	return int64(count), err
}

func applyFilters(query *dbent.ConversationQuery, status, tier, keyword string) {
	if trimmedStatus := strings.TrimSpace(status); trimmedStatus != "" {
		query.Where(entconversation.StatusEQ(trimmedStatus))
	}
	if trimmedTier := strings.TrimSpace(tier); trimmedTier != "" {
		query.Where(entconversation.CustomerTierEQ(trimmedTier))
	}
	if trimmedKeyword := strings.TrimSpace(keyword); trimmedKeyword != "" {
		query.Where(entconversation.Or(
			entconversation.SubjectContainsFold(trimmedKeyword),
			entconversation.CustomerNameContainsFold(trimmedKeyword),
			entconversation.TicketNoContainsFold(trimmedKeyword),
			entconversation.QueueContainsFold(trimmedKeyword),
		))
	}
}

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
)

type Conversation struct{ ent.Schema }

func (Conversation) Mixin() []ent.Mixin { return []ent.Mixin{mixins.Time{}} }

func (Conversation) Fields() []ent.Field {
	return []ent.Field{
		field.String("ticket_no").NotEmpty().Unique(),
		field.String("subject").NotEmpty(),
		field.String("preview").Default(""),
		field.String("channel").Default("Web Chat"),
		field.String("queue").Default("General Queue"),
		field.String("assignee").Default(""),
		field.String("status").Default("waiting"),
		field.String("priority").Default("medium"),
		field.Int("unread").Default(0),
		field.String("sla").Default(""),
		field.String("customer_name").NotEmpty(),
		field.String("customer_company").Default(""),
		field.String("customer_contact").Default(""),
		field.Strings("customer_tags").Default([]string{}),
		field.String("customer_presence").Default("offline"),
		field.String("customer_tier").Default("standard"),
		field.String("last_order").Default(""),
		field.Int("open_tickets").Default(0),
		field.String("satisfaction").Default(""),
		field.Time("last_active_at").Default(time.Now),
		field.Time("closed_at").Optional().Nillable(),
	}
}

func (Conversation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", ConversationMessage.Type),
	}
}

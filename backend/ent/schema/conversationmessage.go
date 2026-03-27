package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ConversationMessage struct{ ent.Schema }

func (ConversationMessage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("conversation_id"),
		field.String("actor").Default(""),
		field.String("actor_type").Default("agent"),
		field.String("message_type").Default("message"),
		field.String("content").Default(""),
		field.Time("created_at").Default(time.Now),
	}
}

func (ConversationMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("conversation", Conversation.Type).Ref("messages").Field("conversation_id").Unique().Required(),
	}
}

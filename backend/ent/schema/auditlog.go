package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type AuditLog struct{ ent.Schema }

func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Default("anonymous"),
		field.String("module").Default("system"),
		field.String("action").Default("view"),
		field.String("method").Default("GET"),
		field.String("path").Default("/"),
		field.Int("status_code").Default(200),
		field.String("ip").Default(""),
		field.String("detail").Default(""),
		field.Time("created_at").Default(time.Now),
	}
}

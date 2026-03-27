package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
)

type SystemConfig struct{ ent.Schema }

func (SystemConfig) Mixin() []ent.Mixin { return []ent.Mixin{mixins.Time{}} }

func (SystemConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Unique(),
		field.String("value").Default(""),
		field.String("category").Default("general"),
		field.String("description").Default(""),
		field.Bool("is_public").Default(false),
	}
}

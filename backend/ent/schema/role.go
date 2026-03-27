package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
)

type Role struct{ ent.Schema }

func (Role) Mixin() []ent.Mixin { return []ent.Mixin{mixins.Time{}} }

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("code").NotEmpty().Unique(),
		field.String("description").Default(""),
		field.JSON("permissions", []string{}).Default([]string{}),
		field.Bool("is_system").Default(false),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{edge.To("users", User.Type)}
}

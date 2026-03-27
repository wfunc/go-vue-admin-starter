package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
)

type User struct{ ent.Schema }

func (User) Mixin() []ent.Mixin { return []ent.Mixin{mixins.Time{}} }

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().Unique(),
		field.String("email").Optional().Nillable().Unique(),
		field.String("password_hash").NotEmpty(),
		field.String("nickname").Default(""),
		field.String("status").Default("active"),
		field.Int("role_id"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("users").Field("role_id").Unique().Required(),
	}
}

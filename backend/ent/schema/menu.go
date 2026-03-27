package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
)

type Menu struct{ ent.Schema }

func (Menu) Mixin() []ent.Mixin { return []ent.Mixin{mixins.Time{}} }

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("name").NotEmpty().Unique(),
		field.String("path").NotEmpty(),
		field.String("component").Default(""),
		field.String("icon").Default(""),
		field.String("menu_type").Default("menu"),
		field.String("permission").Default(""),
		field.Int("sort").Default(0),
		field.Bool("hidden").Default(false),
		field.Int("parent_id").Optional().Nillable(),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Menu.Type),
		edge.From("parent", Menu.Type).Ref("children").Field("parent_id").Unique(),
	}
}

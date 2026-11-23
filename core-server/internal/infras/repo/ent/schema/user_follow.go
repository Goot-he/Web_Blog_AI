package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserFollow struct{ ent.Schema }

func (UserFollow) Annotations() []entschema.Annotation {
	return []entschema.Annotation{entsql.Annotation{Table: "user_follow"}}
}
func (UserFollow) Fields() []ent.Field {
	return []ent.Field{field.Uint64("id"), field.Time("created_at").Immutable().Default(time.Now), field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now), field.Uint64("follower_id"), field.Uint64("following_id"), field.Enum("status").Values("following", "nothing").Default("nothing"), field.Int64("version").Default(0)}
}
func (UserFollow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("follower_id", "following_id").Unique().StorageKey("uk_follow_pair"),
	}
}

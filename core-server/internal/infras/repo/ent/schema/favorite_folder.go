package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type FavoriteFolder struct{ ent.Schema }

func (FavoriteFolder) Annotations() []entschema.Annotation {
	return []entschema.Annotation{entsql.Annotation{Table: "favorite_folder"}}
}
func (FavoriteFolder) Fields() []ent.Field {
	return []ent.Field{field.Uint64("id"), field.Time("created_at").Immutable().Default(time.Now), field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now), field.Uint64("user_id"), field.String("name").MaxLen(50), field.Bool("is_default").Default(false)}
}
func (FavoriteFolder) Indexes() []ent.Index {
	return []ent.Index{index.Fields("user_id", "name").Unique().StorageKey("uk_favorite_folder_user_name")}
}

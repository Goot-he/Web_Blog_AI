package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	entschema "entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InteractionFavor records whether an article is in a user's default favorites collection.
type InteractionFavor struct{ ent.Schema }

func (InteractionFavor) Annotations() []entschema.Annotation {
	return []entschema.Annotation{entsql.Annotation{Table: "interaction_favor"}}
}

func (InteractionFavor) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Uint64("user_id"),
		field.Uint64("article_id"),
		field.Uint64("folder_id"),
		field.Enum("status").Values("favored", "nothing").Default("nothing"),
		field.Int64("version").Default(0),
	}
}

func (InteractionFavor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "article_id").Unique().StorageKey("uk_favor_user_article"),
	}
}

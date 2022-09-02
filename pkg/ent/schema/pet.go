package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Int8("type").NonNegative(),
		field.Time("create_at").Default(time.Now()),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return nil
}

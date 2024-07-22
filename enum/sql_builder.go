package enum

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type Enum struct {
	Value any `json:"value"`
	Enums sqlbuilder.Enums
	Field *sqlbuilder.Field
}

func (b *Enum) Init() {
	b.Field = sqlbuilder.NewField(func(in any) (any, error) { return b.Value, nil }).SetName("enum_column").SetTag("枚举列")
	b.Field.AppendEnum(b.Enums...)
}

func (b Enum) Fields() sqlbuilder.Fields {
	return *sqlbuilder.NewFields(b.Field)
}
func (b *Enum) Apply(initFns ...sqlbuilder.InitFieldFn) *Enum {
	b.Field.Apply(initFns...)
	return b
}

func NewEnum(value any, enums sqlbuilder.Enums) *Enum {
	e := &Enum{
		Value: value,
		Enums: enums,
	}
	e.Init()
	return e
}

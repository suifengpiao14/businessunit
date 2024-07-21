package name

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type Name struct {
	Value string
	Field *sqlbuilder.Field
}

func (name *Name) Init() {
	name.Field = sqlbuilder.NewField(func(in any) (any, error) { return name.Value, nil }).SetName("name").SetTitle("名称")
	name.Field.MergeSchema(sqlbuilder.Schema{
		Required:  false,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
	})
	name.Field.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	name.Field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnWhereLike)
	})

}

func NewName(value string) *Name {
	name := &Name{
		Value: value,
	}
	name.Init()
	return name
}

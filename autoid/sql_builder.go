package autoid

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewAutoIdField(valueFn sqlbuilder.ValueFn) (field *sqlbuilder.Field) {
	field = sqlbuilder.NewField(valueFn)
	field.SetName("id").SetTitle("ID").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_int,
		MaxLength: 64,
		Minimum:   1,
		Primary:   true,
	})

	field.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield)
	})
	field.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // id 不能更新
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})

	field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return field
}

package name

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewNameField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn).SetName("name").SetTitle("名称")
	f.MergeSchema(sqlbuilder.Schema{
		Required:  false,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnWhereLike)
	})
	return f
}

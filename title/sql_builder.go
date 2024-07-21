package title

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewTitleField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn)
	f.SetName("title").SetTitle("标题").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)

	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnWhereLike)
	})
	return f
}

package key

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewKeyField(valueFn sqlbuilder.ValueFn) *sqlbuilder.Field {
	f := sqlbuilder.NewField(valueFn).SetName("key").SetTitle("键")
	f.MergeSchema(sqlbuilder.Schema{
		Required:  false, // insert 场景不一定必须
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)

	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

package key

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionKey(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.SetName("key").SetTitle("键").MergeSchema(sqlbuilder.Schema{
		Required:  false, // insert 场景不一定必须
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
		Primary:   true,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	if field.SceneIsSelect() {
		field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
	}
}

package name

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionName(field *sqlbuilder.Field) {
	field.SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Required:  false,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionName)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionName)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionName).WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect, sqlbuilder.WhereValueFnLike)
}

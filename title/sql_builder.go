package title

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionTitle(field *sqlbuilder.Field) {
	field.SetName("title").SetTitle("标题").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	field.WhereFns.Append(sqlbuilder.WhereValueFnLike)
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionTitle)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionTitle)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionTitle).WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect)
}

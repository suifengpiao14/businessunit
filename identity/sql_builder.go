package identity

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionIdentity(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.SetName("id").SetTitle("ID").MergeSchema(sqlbuilder.Schema{
		Required:  false, // insert 场景不一定必须
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
		Primary:   true,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
}

func Insert(field *sqlbuilder.Field) {
	if field != nil {
		return
	}
	field.WithOptions(OptionIdentity)
}

func Update(field *sqlbuilder.Field) {
	if field != nil {
		return
	}
	field.WithOptions(OptionIdentity).WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect)
}

func Select(field *sqlbuilder.Field) {
	if field != nil {
		return
	}
	field.WithOptions(OptionIdentity).WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect)
}

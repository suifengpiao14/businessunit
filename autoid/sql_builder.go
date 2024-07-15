package autoid

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionAutoID(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.SetName("id").SetTitle("ID").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_int,
		MaxLength: 64,
		Minimum:   1,
		Primary:   true,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionAutoID)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionAutoID).ShieldUpdate(true) // id 不能更新
	field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionAutoID).WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

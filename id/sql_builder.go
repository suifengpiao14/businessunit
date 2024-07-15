package id

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionID(field *sqlbuilder.Field) {
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
	if field == nil {
		return
	}
	field.WithOptions(OptionID)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionID).ShieldUpdate(true) // id 不能更新
	field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionID).WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

package uuid

import (
	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionUUID(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.SetName("uuid").SetTitle("UUID").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		MinLength: 1,
		Primary:   true,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionUUID).ValueFns.InsertAsFirst(func(in any) (any, error) {
		return xid.New().String(), nil
	})
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionUUID).ShieldUpdate(true) // uuid 不能更新
	field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionUUID).WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

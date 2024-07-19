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
	})
	if field.SceneIsInsert() {
		field.SetRequired(true)
		field.ValueFns.InsertAsFirst(func(in any) (any, error) {
			return xid.New().String(), nil
		})
	}

	if field.SceneIsUpdate() {
		field.WithOptions(OptionUUID).ShieldUpdate(true) // uuid 不能更新
		field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnEmpty2Nil, sqlbuilder.ValueFnForward)
	}

	if field.SceneIsSelect() {
		field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
	}
}

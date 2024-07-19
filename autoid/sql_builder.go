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
	})
	if field.SceneIsInsert() {
		field.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	} else if field.SceneIsUpdate() {
		field.WithOptions(OptionAutoID).ShieldUpdate(true) // id 不能更新
		field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
	} else if field.SceneIsSelect() {
		field.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	}
}

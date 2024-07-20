package ownerid

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionOwnerID(f *sqlbuilder.Field) {
	f.SetName("ownerId").SetTitle("所有者").MergeSchema(sqlbuilder.Schema{
		Comment:      "对象标识,缺失时记录无意义",
		Type:         sqlbuilder.Schema_Type_string,
		MaxLength:    64,
		MinLength:    1,
		Minimum:      1,
		ShieldUpdate: true, // 所有者不可跟新
	})
	f.SceneInsert(func(f *sqlbuilder.Field) {
		f.SetRequired(true)
	})

	f.SceneUpdate(func(f *sqlbuilder.Field) {
		f.ShieldUpdate(true)
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})

	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.ValueFns.AppendIfNotFirst(sqlbuilder.ValueFnEmpty2Nil)
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
}

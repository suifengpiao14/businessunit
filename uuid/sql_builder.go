package uuid

import (
	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

func NewUuidField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn).SetName("uuid").SetTitle("UUID")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		MinLength: 1,
	})
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetRequired(true)
		f.ValueFns.InsertAsFirst(func(in any) (any, error) {
			return xid.New().String(), nil
		})
	})
	f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // uuid 不能更新
		f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnEmpty2Nil, sqlbuilder.ValueFnForward)
	})

	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
	})

	return f
}

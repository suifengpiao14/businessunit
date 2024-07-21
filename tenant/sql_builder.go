package tenant

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewTenantField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn).SetName("ternatId").SetTitle("租户ID")
	f.MergeSchema(sqlbuilder.Schema{
		Required:  true,
		MinLength: 1,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Maximum:   sqlbuilder.UnsinedInt_maximum_bigint,
		Minimum:   1,
	})
	f.SceneUpdate(func(f *sqlbuilder.Field) {
		f.ShieldUpdate(true) // 不可更新
	})
	f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward) // update,select 都必须为条件
	return f
}

package tenant

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionTenant(f *sqlbuilder.Field) {
	f.SetName("ternatId").SetTitle("租户ID").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		MinLength: 1,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Maximum:   sqlbuilder.UnsinedInt_maximum_bigint,
		Minimum:   1,
	})

}

func Insert(f *sqlbuilder.Field) {
	if f == nil {
		return
	}
	f.WithOptions(OptionTenant)
}

func Update(f *sqlbuilder.Field) {
	if f == nil {
		return
	}
	f.WithOptions(OptionTenant)
	f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward) // update,select 都必须为条件
}

func Select(f *sqlbuilder.Field) {
	if f == nil {
		return
	}
	f.WithOptions(OptionTenant)
	f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward) // update,select 都必须为条件
}

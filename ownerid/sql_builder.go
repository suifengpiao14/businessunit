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
}

func Insert(f *sqlbuilder.Field) {
	f.WithOptions(OptionOwnerID).SetRequired(true) // 新增时不能为空
}

func Update(f *sqlbuilder.Field) {
	f.WithOptions(OptionOwnerID).ShieldUpdate(true) // 不可更新
	f.ValueFns.AppendIfNotFirst(sqlbuilder.ValueFnEmpty2Nil)
	f.WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect)

}
func Select(f *sqlbuilder.Field) {
	f.WithOptions(OptionOwnerID)
	f.ValueFns.AppendIfNotFirst(sqlbuilder.ValueFnEmpty2Nil)
	f.WhereFns.InsertAsSecond(sqlbuilder.WhereValueFnDirect)

}

package phone

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionPhone(f *sqlbuilder.Field) {
	f.SetName("phone").SetTitle("手机号").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 15,
		RegExp:    `^1[3-9]\d{9}$`, // 中国大陆手机号正则表达式
	})
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionPhone)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionPhone)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionPhone).WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnDirect)
}

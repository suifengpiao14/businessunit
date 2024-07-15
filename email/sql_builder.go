package email

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionEmail(f *sqlbuilder.Field) {
	f.SetName("email").SetTitle("邮箱").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 32,
		MinLength: 5,
		RegExp:    `([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}`, // 邮箱验证表达式
	})
}

func Insert(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionEmail)
}

func Update(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionEmail).WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

func Select(field *sqlbuilder.Field) {
	if field == nil {
		return
	}
	field.WithOptions(OptionEmail).WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
}

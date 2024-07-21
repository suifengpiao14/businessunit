package email

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewEmailField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn).SetName("email").SetTitle("邮箱")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 32,
		MinLength: 5,
		RegExp:    `([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}`, // 邮箱验证表达式
	})
	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward)
	})
	return f
}

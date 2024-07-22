package phone

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type Phone struct {
	Value any `json:"phone"`
	Field *sqlbuilder.Field
}

var DefaultInitFieldFn sqlbuilder.InitFieldFn = func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
	f.SetName("phone").SetTitle("手机号").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 15,
		RegExp:    `^1[3-9]\d{9}$`, // 中国大陆手机号正则表达式
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
}

func (p *Phone) Apply(initFns ...sqlbuilder.InitFieldFn) *Phone {
	p.Field.Apply(initFns...)
	return p
}

func (p Phone) Fields() sqlbuilder.Fields {
	return *sqlbuilder.NewFields(p.Field)
}

func NewPhone(value any) *Phone {
	p := &Phone{
		Value: value,
		Field: sqlbuilder.NewField(func(in any) (any, error) { return value, nil }),
	}
	p.Apply(DefaultInitFieldFn)
	return p
}

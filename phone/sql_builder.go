package phone

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type PhoneField struct {
	sqlbuilder.Field
}

func (f PhoneField) GetPhoneField() PhoneField {
	return f
}

func (f PhoneField) SetName(name string) PhoneField {
	f.Field.SetName(name)
	return f
}
func (f PhoneField) SetTitle(title string) PhoneField {
	f.Field.SetTitle(title)
	return f
}

// AppendWhereFn 添加Where条件，方便连续书写
func (f PhoneField) AppendWhereFn(fns ...sqlbuilder.WhereValueFn) PhoneField {
	f.Field.WhereFns.Append(fns...)
	return f
}

var PhoneFieldSchema = sqlbuilder.Schema{
	Type:      sqlbuilder.Schema_Type_string,
	MaxLength: 15,
	RegExp:    `^1[3-9]\d{9}$`, // 中国大陆手机号正则表达式
}

func NewPhoneField(valueFn sqlbuilder.ValueFn) (field PhoneField) {
	field = PhoneField{
		Field: *sqlbuilder.NewField(valueFn).SetName("phone").SetTitle("手机号").MergeSchema(PhoneFieldSchema),
	}
	return field
}

type PhoneI interface {
	GetPhoneField() PhoneField // 使用每个包下重命名的类型，具有区分度
}

func _DataFn(phoneI PhoneI) sqlbuilder.DataFn {
	field := phoneI.GetPhoneField()
	return field.Data
}

func _WhereFn(phoneI PhoneI) sqlbuilder.WhereFn {
	return phoneI.GetPhoneField().Where
}

func Insert(phoneI PhoneI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(phoneI))
}

func Update(phoneI PhoneI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(phoneI))
}

func First(phoneI PhoneI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(phoneI))
}

func List(phoneI PhoneI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(phoneI))
}

func Total(phoneI PhoneI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(phoneI))
}

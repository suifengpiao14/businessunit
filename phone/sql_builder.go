package phone

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/suifengpiao14/sqlbuilder"
)

// 定义一个全局的验证器
var validate = validator.New()

func init() {
	// 注册自定义验证函数
	validate.RegisterValidation("phone", validateMobile)
}

type PhoneField struct {
	sqlbuilder.Field
}

func (f PhoneField) GetPhoneField() PhoneField {
	return f
}

type PhoneI interface {
	GetPhoneField() PhoneField // 使用每个包下重命名的类型，具有区分度
}

func validatePhone(phone string) (err error) {
	err = validate.Var(phone, "phone")
	return err
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

// 自定义验证函数，使用正则表达式验证手机号格式
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 中国大陆手机号正则表达式
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(mobile)
}

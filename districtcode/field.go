package districtcode

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionCode(f *sqlbuilder.Field) {
	f.SetName("code").SetTitle("行政区代码").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 16, //统计局统一是使用12位，如 659008103505
		MinLength: 2,  // 只使用省时，为2位
	})
}

func OptionName(f *sqlbuilder.Field) {
	f.SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64, //线上统计最长55个汉字，使用64符合大部分场景
		MinLength: 2,  // 只使用省时，为2个汉字
	})
}

func Options(codeFiled *sqlbuilder.Field, nameField *sqlbuilder.Field) {
	codeFiled.WithOptions(OptionCode)
	nameField.WithOptions(OptionName)
}

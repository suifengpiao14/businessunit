package districtcode

import (
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/idtitle"
	"github.com/suifengpiao14/sqlbuilder"
)

func NewDistrictCodeField(valueFn sqlbuilder.ValueFn) identity.IdentityField {
	field := *identity.NewIdentityField(valueFn)
	field.SetName("code").SetTitle("行政区代码").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 16, //统计局统一是使用12位，如 659008103505
		MinLength: 2,  // 只使用省时，为2位
	})

	return field
}

func NewDistrictNameField(valueFn sqlbuilder.ValueFn) sqlbuilder.Field {
	field := idtitle.NewTitleField(valueFn)
	field.SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64, //线上统计最长55个汉字，使用64符合大部分场景
		MinLength: 2,  // 只使用省时，为2个汉字
	})
	return *field
}

type DistrictField struct {
	*idtitle.IdTitleField
}

func NewDistrictField(codeValueFn sqlbuilder.ValueFn, nameValueFn sqlbuilder.ValueFn) DistrictField {
	id, title := NewDistrictCodeField(codeValueFn), NewDistrictNameField(nameValueFn)
	field := &idtitle.IdTitleField{
		ID:    &id,
		Title: &title,
	}
	return DistrictField{
		field,
	}
}

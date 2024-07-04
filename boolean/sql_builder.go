package boolean

import (
	"strings"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/sqlbuilder"
)

type BooleanField struct {
	sqlbuilder.Field
	TrueFalseTitleFn func() (trueTitle enum.EnumTitle, falseTitle enum.EnumTitle)
}

func (f BooleanField) GetBooleanField() BooleanField {
	return f
}

func (f BooleanField) GetTrueFalseTitle() (trueTitle enum.EnumTitle, falseTitle enum.EnumTitle) {
	return f.TrueFalseTitleFn()
}

func (f BooleanField) IsTrue() (isTrue bool) {
	if f.ValueFn == nil {
		return false
	}
	val, err := f.ValueFn(nil)
	if err != nil {
		return false
	}
	trueTitle, _ := f.GetTrueFalseTitle()
	isTrue = strings.EqualFold(trueTitle.Key, cast.ToString(val))
	return isTrue
}

type BooleanI interface {
	GetBooleanField() BooleanField
	GetTrueFalseTitle() (trueTitle enum.EnumTitle, falseTitle enum.EnumTitle)
}

// Copy 生成副本，此处实际类型可能已经发生变化，只是复制了BooleanI 接口内容
func Copy(booleanI BooleanI) (newBooleanI BooleanI) {
	booleanField := booleanI.GetBooleanField()
	newBooleanI = BooleanField{
		Field: sqlbuilder.Field{
			Name:    booleanField.Name,
			ValueFn: booleanField.ValueFn,
		},
		TrueFalseTitleFn: booleanI.GetTrueFalseTitle,
	}
	return newBooleanI
}

// Switch  生成值反转的实体
func Switch(booleanI BooleanI) (reversed BooleanI) {
	booleanField := booleanI.GetBooleanField()
	trueTitle, falseTitle := booleanI.GetTrueFalseTitle()
	reversed = BooleanField{
		Field: sqlbuilder.Field{
			Name: booleanField.Name,
			ValueFn: func(in any) (value any, err error) {
				val, err := booleanField.ValueFn(nil)
				if err != nil {
					return nil, err
				}
				str := cast.ToString(val)
				if trueTitle.IsSame(str) { // 原值为true ，返回 false 值
					return falseTitle.Key, nil
				}
				return trueTitle.Key, nil // 原值为false ，返回 true 值

			},
		},
		TrueFalseTitleFn: booleanI.GetTrueFalseTitle,
	}
	return reversed
}

func booleanI2EnumField(booleanI BooleanI) enum.EnumField {
	trueTitle, falseTitle := booleanI.GetTrueFalseTitle()
	enumFiled := enum.EnumField{
		Field: sqlbuilder.Field(booleanI.GetBooleanField().Field),
		EnumTitles: enum.EnumTitles{
			trueTitle, falseTitle,
		},
	}
	return enumFiled
}

func Insert(booleanI BooleanI) sqlbuilder.InsertParam {
	if booleanI == nil {
		return sqlbuilder.InsertParam{}
	}
	return enum.Insert(booleanI2EnumField(booleanI))
}

func Update(booleanI BooleanI) sqlbuilder.UpdateParam {
	if booleanI == nil {
		return sqlbuilder.UpdateParam{}
	}
	return enum.Update(booleanI2EnumField(booleanI))
}

func First(booleanI BooleanI) sqlbuilder.FirstParam {
	if booleanI == nil {
		return sqlbuilder.FirstParam{}
	}
	return enum.First(booleanI2EnumField(booleanI))
}

func List(booleanI BooleanI) sqlbuilder.ListParam {
	if booleanI == nil {
		return sqlbuilder.ListParam{}
	}
	return enum.List(booleanI2EnumField(booleanI))
}

func Total(booleanI BooleanI) sqlbuilder.TotalParam {
	if booleanI == nil {
		return sqlbuilder.TotalParam{}
	}
	return enum.Total(booleanI2EnumField(booleanI))
}

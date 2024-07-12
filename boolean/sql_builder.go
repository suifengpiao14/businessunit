package boolean

import (
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/sqlbuilder"
)

type TrueFalseTitleFn func() (trueTitle sqlbuilder.Enum, falseTitle sqlbuilder.Enum)

type BooleanField struct {
	sqlbuilder.Field
	TrueFalseTitleFn TrueFalseTitleFn
}

func NewBooleanField(valueFn sqlbuilder.ValueFn, trueFalseTitleFn TrueFalseTitleFn) BooleanField {
	trueTitle, falseTitle := trueFalseTitleFn()
	schema := sqlbuilder.Schema{
		Enums: sqlbuilder.Enums{
			{
				Title: trueTitle.Title,
				Key:   trueTitle.Key,
			},
			{
				Title: falseTitle.Title,
				Key:   falseTitle.Key,
			},
		},
	}
	return BooleanField{
		Field:            *sqlbuilder.NewField(valueFn).MergeSchema(schema),
		TrueFalseTitleFn: trueFalseTitleFn,
	}
}

func (f BooleanField) GetBooleanField() BooleanField {
	return f
}

func (f BooleanField) GetTrueFalseTitle() (trueTitle sqlbuilder.Enum, falseTitle sqlbuilder.Enum) {
	return f.TrueFalseTitleFn()
}

func (f BooleanField) AppendWhereFn(fns ...sqlbuilder.ValueFn) BooleanField {
	f.Field.AppendWhereFn(fns...)
	return f
}

func (f BooleanField) IsTrue() (isTrue bool) {
	if f.ValueFns == nil {
		return false
	}
	val, err := f.GetValue(nil)
	if err != nil {
		return false
	}
	trueTitle, _ := f.GetTrueFalseTitle()
	isTrue = trueTitle.IsEqual(val)
	return isTrue
}

type BooleanI interface {
	GetBooleanField() BooleanField
	GetTrueFalseTitle() (trueTitle sqlbuilder.Enum, falseTitle sqlbuilder.Enum)
}

// Copy 生成副本，此处实际类型可能已经发生变化，只是复制了BooleanI 接口内容
func Copy(booleanI BooleanI) (newBooleanI BooleanI) {
	booleanField := booleanI.GetBooleanField()
	newBooleanI = BooleanField{
		Field: sqlbuilder.Field{
			ValueFns: booleanField.ValueFns,
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
			ValueFns: sqlbuilder.ValueFns{
				func(in any) (value any, err error) {
					val, err := booleanField.GetValue(nil)
					if err != nil {
						return nil, err
					}
					if trueTitle.IsEqual(val) { // 原值为true ，返回 false 值
						return falseTitle.Key, nil
					}
					return trueTitle.Key, nil // 原值为false ，返回 true 值

				},
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
		EnumTitles: sqlbuilder.Enums{
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

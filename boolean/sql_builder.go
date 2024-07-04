package boolean

import (
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

type BooleanI interface {
	GetBooleanField() BooleanField
	GetTrueFalseTitle() (trueTitle enum.EnumTitle, falseTitle enum.EnumTitle)
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

package boolean

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionBooleanFn(trueTitle sqlbuilder.Enum, falseTitle sqlbuilder.Enum) sqlbuilder.OptionFn {
	return func(field *sqlbuilder.Field) {
		field.SetName("boolean").SetTitle("布尔值")
		trueTitle.Tag = sqlbuilder.Enum_tag_true
		falseTitle.Tag = sqlbuilder.Enum_tag_false
		field.MergeSchema(sqlbuilder.Schema{
			Enums: sqlbuilder.Enums{
				trueTitle,
				falseTitle,
			},
		})
	}
}

// Switch  将值反转
func Switch(f *sqlbuilder.Field) {
	f.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := f.Schema.Enums
		for _, enum := range enums {
			if !enum.IsEqual(in) {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found reversed enum key ;current:%v", in)
		return nil, err
	})
}

// TrunOff  改成false
func TrunOff(f *sqlbuilder.Field) {
	f.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := f.Schema.Enums
		for _, enum := range enums {
			if enum.Tag == sqlbuilder.Enum_tag_false {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found fase enum key enums:%s", enums.String())
		return nil, err
	})
}

// TrunOn  改成true
func TrunOn(f *sqlbuilder.Field) {
	f.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := f.Schema.Enums
		for _, enum := range enums {
			if enum.Tag == sqlbuilder.Enum_tag_true {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found true enum key enums:%s", enums.String())
		return nil, err
	})
}

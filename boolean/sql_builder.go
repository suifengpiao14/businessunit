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
		err := errors.Errorf("not foun reversed enum key ;current:%v", in)
		return nil, err
	})

}

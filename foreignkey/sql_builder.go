package foreignkey

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionForeignkey(f *sqlbuilder.Field, redundantFields ...sqlbuilder.Field) {
	if f.SceneIsInsert() && len(redundantFields) > 0 {
		f.ValueFns.InsertAsSecond(func(in any) (any, error) {
			val, err := f.GetValue()
			if err != nil {
				return nil, err
			}
			m := map[string]any{}
			for _, redundantField := range redundantFields {
				redundantField.ValueFns.InsertAsFirst(func(in any) (any, error) { return val, nil })
				redundantFiledValue, err := redundantField.GetValue()
				if err != nil {
					return nil, err
				}
				if !sqlbuilder.IsNil(redundantFiledValue) {
					m[redundantField.DBName()] = redundantFiledValue
				}
			}
			return m, nil
		})

	}

}

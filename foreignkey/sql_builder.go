package foreignkey

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionForeignkey(f *sqlbuilder.Field, redundantFields ...sqlbuilder.Field) {
	if len(redundantFields) > 0 {
		return
	}
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFn{
			Layer: sqlbuilder.Value_Layer_ApiFormat,
			Fn: func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
				val, err := f.GetValue(sqlbuilder.Layer_all)
				if err != nil {
					return nil, err
				}
				m := map[string]any{}
				for _, redundantField := range redundantFields {
					redundantField.ValueFns.Append(sqlbuilder.ValueFn{
						Layer: sqlbuilder.Value_Layer_SetValue,
						Fn:    func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return val, nil },
					})
					redundantFiledValue, err := redundantField.GetValue(sqlbuilder.Layer_all)
					if err != nil {
						return nil, err
					}
					if !sqlbuilder.IsNil(redundantFiledValue) {
						m[redundantField.DBName()] = redundantFiledValue
					}
				}
				return m, nil
			},
		})
	})

}

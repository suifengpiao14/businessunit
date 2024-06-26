package unique

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type UniqueField sqlbuilder.Fields

func (f UniqueField) GetUniqueFields() UniqueField {
	return f
}

type UniqueI interface {
	GetUniqueFields() UniqueField
}

func _DataFn(uniqueI UniqueI) sqlbuilder.DataFn {
	return func() (any, error) {
		fields := uniqueI.GetUniqueFields()
		m := map[string]any{}
		for _, field := range fields {
			val, err := field.Value(nil)
			if err != nil {
				return nil, err
			}
			m[field.Name] = val
		}
		return m, nil
	}
}

func _whereFn(uniqueI UniqueI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		fields := uniqueI.GetUniqueFields()
		expressions = make([]goqu.Expression, 0)
		for _, field := range fields {
			val, err := field.Value(nil)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, goqu.C(field.Name).Eq(val))
		}
		return expressions, nil
	}
}

func Exists(uniqueI UniqueI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_whereFn(uniqueI))
}

func Insert(uniqueI UniqueI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(uniqueI))
}

func Update(uniqueI UniqueI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(uniqueI))
}

func First(uniqueI UniqueI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(uniqueI UniqueI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(uniqueI UniqueI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

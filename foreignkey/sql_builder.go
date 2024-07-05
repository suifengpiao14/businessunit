package foreignkey

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type ForeignkeyField sqlbuilder.Field

func (f ForeignkeyField) GetForeignkeyField() ForeignkeyField {
	return f
}

type ForeignkeyI interface {
	GetForeignkeyField() ForeignkeyField
	RedundantFields() sqlbuilder.Fields
}

func _DataFn(foreignkeyI ForeignkeyI) sqlbuilder.DataFn {
	return func() (any, error) {
		col := foreignkeyI.GetForeignkeyField()
		if col.ValueFns == nil {
			return nil, nil
		}
		val, err := sqlbuilder.Field(col).GetValue(nil)
		if err != nil {
			return nil, err
		}
		m := map[string]any{}
		redundantFields := foreignkeyI.RedundantFields()
		for _, redundantField := range redundantFields {
			redundantFiledValue, err := redundantField.GetValue(val)
			if err != nil {
				return nil, err
			}
			if !sqlbuilder.IsNil(redundantFiledValue) {
				m[redundantField.Name] = redundantFiledValue
			}
		}
		return m, nil
	}
}

func Insert(foreignkeyI ForeignkeyI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(foreignkeyI))
}

func Update(foreignkeyI ForeignkeyI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(foreignkeyI))
}

func First(foreignkeyI ForeignkeyI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(foreignkeyI ForeignkeyI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(foreignkeyI ForeignkeyI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

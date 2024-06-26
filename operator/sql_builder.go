package operator

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type OperatorField struct {
	OperatorID   sqlbuilder.Field
	OperatorName sqlbuilder.Field
}

type OperatorI interface {
	GetOperatorField() OperatorField
}

func _OperatorFn(operatorI OperatorI) sqlbuilder.DataFn {
	return func() (any, error) {
		col := operatorI.GetOperatorField()
		m := map[string]any{}
		if col.OperatorID.Name != "" {
			val, err := col.OperatorID.Value(nil)
			if err != nil {
				return nil, err
			}
			m[col.OperatorID.Name] = val
		}
		if col.OperatorName.Name != "" {
			val, err := col.OperatorName.Value(nil)
			if err != nil {
				return nil, err
			}
			m[col.OperatorName.Name] = val
		}
		return m, nil
	}
}

func Insert(operatorI OperatorI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_OperatorFn(operatorI))
}

func Update(operatorI OperatorI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_OperatorFn(operatorI))
}

func First(operatorI OperatorI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(operatorI OperatorI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(operatorI OperatorI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

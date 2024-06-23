package operator

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type OperatorField struct {
	OperatorID   sqlbuilder.Column
	OperatorName sqlbuilder.Column
}

type OperatorI interface {
	GetOperatorField() OperatorField
}

func _OperatorFn(operatorI OperatorI) sqlbuilder.DataFn {
	col := operatorI.GetOperatorField()
	m := map[string]any{}
	if col.OperatorID.Name != "" {
		m[col.OperatorID.Name] = col.OperatorID.Value(nil)
	}
	if col.OperatorName.Name != "" {
		m[col.OperatorName.Name] = col.OperatorName.Value(nil)
	}
	return func() (any, error) {
		return m, nil
	}
}

func Insert(operatorI OperatorI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_OperatorFn(operatorI))
}

func Update(operatorI OperatorI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_OperatorFn(operatorI))
}

var SorftDelete = Update

func First(operatorI OperatorI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(operatorI OperatorI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(operatorI OperatorI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

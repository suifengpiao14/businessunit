package operator

import (
	"github.com/suifengpiao14/businessunit"
	"github.com/suifengpiao14/sqlbuilder"
)

type Operator[T int | string] struct {
	OperatorId   T      `json:"opreatorId"`
	OperatorName string `json:"opreatorName"`
}

func (o Operator[T]) Feilds() *OperatorFields {

	return &OperatorFields{
		businessunit.NewKeyFieldField(o.OperatorId).SetName("opreatorId").SetTitle("操作人ID"),
		businessunit.NewNameField(o.OperatorName).SetName("operatorName").SetTitle("操作人名称"),
	}
}

type OperatorFields struct {
	OperatorIdField   *sqlbuilder.Field
	OperatorNameField *sqlbuilder.Field
}

func (o OperatorFields) Fields() (fs sqlbuilder.Fields) {
	fs = sqlbuilder.Fields{
		o.OperatorIdField,
		o.OperatorNameField,
	}
	return fs
}

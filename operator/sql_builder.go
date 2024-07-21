package operator

import (
	"github.com/suifengpiao14/businessunit/key"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/sqlbuilder"
)

type Operator struct {
	OperatorId        any    `json:"opreatorId"`
	OperatorName      string `json:"opreatorName"`
	operatorIdField   *sqlbuilder.Field
	operatorNameField *sqlbuilder.Field
}

func (o Operator) Fields() (fs sqlbuilder.Fields) {
	fs = sqlbuilder.Fields{
		o.operatorIdField,
		o.operatorNameField,
	}
	return fs
}

func (o Operator) InitField() {
	o.operatorIdField = key.NewKeyField(func(in any) (any, error) { return o.OperatorId, nil }).SetName("opreatorId").SetTitle("操作人ID")
	o.operatorNameField = name.NewNameField(func(in any) (any, error) { return o.OperatorName, nil }).SetName("operatorName").SetTitle("操作人名称")
}

package operator

import (
	"github.com/suifengpiao14/businessunit/key"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/sqlbuilder"
)

type Operator struct {
	OperatorId      any    `json:"opreatorId"`
	OperatorName    string `json:"opreatorName"`
	operatorIdField *sqlbuilder.Field
	operatorName    *name.Name
}

func (o Operator) Fields() (fs sqlbuilder.Fields) {
	fs = sqlbuilder.Fields{
		o.operatorIdField,
		o.operatorName.Field,
	}
	return fs
}

func (o Operator) InitField() {
	o.operatorIdField = key.NewKeyField(func(in any) (any, error) { return o.OperatorId, nil }).SetName("opreatorId").SetTitle("操作人ID")
	o.operatorName = name.NewName(o.OperatorName)
	o.operatorName.Field.SetName("operatorName").SetTitle("操作人名称")
}

package operator

import (
	"github.com/suifengpiao14/businessunit/keytitle"
	"github.com/suifengpiao14/sqlbuilder"
)

func Insert(operatorIdField *sqlbuilder.Field, operatorNameField *sqlbuilder.Field) {
	keytitle.Insert(operatorIdField, operatorNameField)
}

func Update(operatorIdField *sqlbuilder.Field, operatorNameField *sqlbuilder.Field) {
	keytitle.Insert(operatorIdField, operatorNameField)
}

func Select(operatorIdField *sqlbuilder.Field, operatorNameField *sqlbuilder.Field) {
	keytitle.Insert(operatorIdField, operatorNameField)
}

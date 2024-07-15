package keytitle

import (
	"github.com/suifengpiao14/businessunit/key"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

func Insert(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(key.Insert)
	}
	if titleField != nil {
		titleField.WithOptions(title.Insert)
	}
}

func Update(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(key.Update)
	}
	if titleField != nil {
		titleField.WithOptions(title.Update)
	}
}

func Select(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(key.Update)
	}
	if titleField != nil {
		titleField.WithOptions(title.Update)
	}
}

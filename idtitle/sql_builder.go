package idtitle

import (
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

func Insert(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(identity.Insert)
	}
	if titleField != nil {
		titleField.WithOptions(title.Insert)
	}
}

func Update(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(identity.Update)
	}
	if titleField != nil {
		titleField.WithOptions(title.Update)
	}
}

func First(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(identity.First)
	}
	if titleField != nil {
		titleField.WithOptions(title.First)
	}
}

func List(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(identity.List)
	}
	if titleField != nil {
		titleField.WithOptions(title.List)
	}
}

func Total(idField *sqlbuilder.Field, titleField *sqlbuilder.Field) {
	if idField != nil {
		idField.WithOptions(identity.Total)
	}
	if titleField != nil {
		titleField.WithOptions(title.Total)
	}
}

package keytitle

import (
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/sqlbuilder"
)

type IdName struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IdField   *identity.Identifier
	NameField *name.Name
}

func (idName IdName) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		idName.IdField.Field,
		idName.NameField.Field,
	}
	return fs
}

func (idName *IdName) Init() {
	idName.IdField = identity.NewIdentifier(idName.Id)
	idName.NameField = name.NewName(idName.Name)
}

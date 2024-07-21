package keytitle

import (
	"github.com/suifengpiao14/businessunit/key"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

type KeyTitle struct {
	Key        string `json:"key"`
	Title      string `json:"title"`
	KeyField   *sqlbuilder.Field
	TitleField *sqlbuilder.Field
}

func (kt KeyTitle) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		kt.KeyField,
		kt.TitleField,
	}
	return fs
}

func (kt KeyTitle) InitField() {
	kt.KeyField = key.NewKeyField(func(in any) (any, error) { return kt.Key, nil })
	kt.TitleField = title.NewTitleField(func(in any) (any, error) { return kt.Title, nil })
}

package businessunit

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewNickname(nickname string) *sqlbuilder.Field {
	f := NewNameField(nickname).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("nickname").SetTitle("昵称")
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(
			sqlbuilder.ValueFnEmpty2Nil,
			sqlbuilder.ValueFnWhereLike,
		)
	})
	return f
}

type KeyTitleField struct {
	KeyField   *sqlbuilder.Field
	TitleField *sqlbuilder.Field
}

func (kt KeyTitleField) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		kt.KeyField,
		kt.TitleField,
	}
	return fs
}

func NewKeyTitleField(key any, title string) *KeyTitleField {
	return &KeyTitleField{
		NewKeyField(key),
		NewTitleField(title),
	}
}

type IdNameFields struct {
	IdField   *sqlbuilder.Field
	NameField *sqlbuilder.Field
}

func (idName IdNameFields) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		idName.IdField,
		idName.NameField,
	}
	return fs
}

func NewIdNameFields(id any, name string) *IdNameFields {
	return &IdNameFields{
		IdField:   NewIdentifierField(id),
		NameField: NewNameField(name),
	}

}

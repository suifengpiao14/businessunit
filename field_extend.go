package businessunit

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func NewNickname(nickname string) *sqlbuilder.Field {
	f := NewNameField(nickname).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("nickname").SetTitle("昵称")
	})
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(
			sqlbuilder.ValueFnEmpty2Nil,
			sqlbuilder.ValueFnWhereLike,
		)
	})
	return f
}

type KeyTitleFields[keyT int | string] struct {
	Key   sqlbuilder.FieldFn[keyT]
	Title sqlbuilder.FieldFn[string]
}

func (KeyTitleFields[T]) Fields() KeyTitleFields[T] {

	return KeyTitleFields[T]{
		Key:   NewKeyField[T],
		Title: NewTitleField,
	}
}

type IdNameFields[idT int | string] struct {
	IdField   sqlbuilder.FieldFn[idT]
	NameField sqlbuilder.FieldFn[string]
}

func (IdNameFields[T]) Builder() IdNameFields[T] {
	return IdNameFields[T]{
		IdField: func(value T) *sqlbuilder.Field {
			return NewIdentifierField(value)
		},
		NameField: func(value string) *sqlbuilder.Field {
			return NewNameField(value)
		},
	}
}

type CUTimeFields struct {
	CreatedAt sqlbuilder.FieldFn[string]
	UpdatedAt sqlbuilder.FieldFn[string]
}

func (CUTimeFields) Builder() CUTimeFields {
	return CUTimeFields{
		CreatedAt: func(value string) *sqlbuilder.Field {
			return NewCreatedAtField()
		},
		UpdatedAt: func(value string) *sqlbuilder.Field {
			return NewUpdatedAtField()
		},
	}
}

type CUDTimeFields struct {
	CreatedAt sqlbuilder.FieldFn[string]
	UpdatedAt sqlbuilder.FieldFn[string]
	DeletedAt sqlbuilder.FieldFn[string]
}

func (CUDTimeFields) Builder() CUDTimeFields {
	cuTime := new(CUTimeFields).Builder()
	return CUDTimeFields{
		CreatedAt: cuTime.CreatedAt,
		UpdatedAt: cuTime.UpdatedAt,
		DeletedAt: func(value string) *sqlbuilder.Field {
			return NewDeletedAtField()
		},
	}
}

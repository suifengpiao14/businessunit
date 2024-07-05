package identity

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type IdentityField sqlbuilder.Field

func (f IdentityField) GetIdentityField() IdentityField {
	return f
}

type IdentityI interface {
	GetIdentityField() IdentityField
}

func _DataFn(identityI IdentityI) sqlbuilder.DataFn {
	return func() (any, error) {
		field := identityI.GetIdentityField()
		if field.ValueFns == nil {
			return nil, nil
		}
		val, err := sqlbuilder.Field(field).GetValue(nil)
		if err != nil {
			return nil, err
		}

		m := map[string]any{
			field.Name: val,
		}
		return m, nil
	}
}

func _WhereFn(identityI IdentityI) sqlbuilder.WhereFn {
	field := identityI.GetIdentityField()
	return sqlbuilder.Field(field).Where
}

func Insert(identityI IdentityI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(identityI))
}

func Update(identityI IdentityI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(identityI)).AppendWhere(_WhereFn(identityI))
}

func First(identityI IdentityI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(identityI))
}

func List(identityI IdentityI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(identityI))
}

func Total(identityI IdentityI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(identityI))
}

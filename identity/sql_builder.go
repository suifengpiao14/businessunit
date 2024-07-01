package identity

import (
	"github.com/doug-martin/goqu/v9"
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
		val, err := field.Value(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, nil
		}

		m := map[string]any{
			field.Name: val,
		}
		return m, nil
	}
}

func _WhereFn(identityI IdentityI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := identityI.GetIdentityField()
		expressions = make([]goqu.Expression, 0)
		val, err := field.WhereValue(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, nil
		}
		if ex, ok := sqlbuilder.TryConvert2Expressions(val); ok {
			return ex, nil
		}
		expressions = append(expressions, goqu.Ex{field.Name: val}) // 支持数组
		return expressions, nil
	}
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

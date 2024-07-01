package tenant

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

type TenantField sqlbuilder.Field

func (f TenantField) GetTenantField() TenantField {
	return f
}

type TenantI interface {
	GetTenantField() TenantField // 获取多租户标记的字段名及值
}

func _DataFn(tenantI TenantI) sqlbuilder.DataFn {
	field := tenantI.GetTenantField()
	return func() (any, error) {
		if field.Value == nil {
			return nil, nil
		}
		m := map[string]any{}
		val, err := field.Value(nil)
		if err != nil {
			return nil, err
		}
		m[field.Name] = val
		return m, nil
	}
}

func _whereFn(uniqueI TenantI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := uniqueI.GetTenantField()
		expressions = make([]goqu.Expression, 0)
		val, err := field.WhereValue(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, err
		}
		expressions = append(expressions, goqu.C(field.Name).Eq(val))
		return expressions, nil
	}
}

func Insert(tenant TenantI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(tenant))
}

func Update(tenant TenantI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(_whereFn(tenant))
}

func First(tenant TenantI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_whereFn(tenant))
}

func List(tenant TenantI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_whereFn(tenant))
}

func Total(tenant TenantI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_whereFn(tenant))
}

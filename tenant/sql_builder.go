package tenant

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type TenantField struct {
	sqlbuilder.Field
}

func (f TenantField) GetTenantField() TenantField {
	return f
}
func (f TenantField) IsEqual(o TenantField) bool {
	return sqlbuilder.Field(f.Field).IsEqual(o.Field)
}

type TenantI interface {
	GetTenantField() TenantField // 获取多租户标记的字段名及值
}

func _DataFn(tenantI TenantI) sqlbuilder.DataFn {
	return tenantI.GetTenantField().Data
}

func WhereFn(uniqueI TenantI) sqlbuilder.WhereFn {
	return uniqueI.GetTenantField().Where
}

func Insert(tenant TenantI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(tenant))
}

func Update(tenant TenantI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(WhereFn(tenant))
}

func First(tenant TenantI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(WhereFn(tenant))
}

func List(tenant TenantI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(WhereFn(tenant))
}

func Total(tenant TenantI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(WhereFn(tenant))
}

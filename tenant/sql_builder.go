package tenant

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

type TenantI interface {
	GetTenant() (filedName string, value any) // 获取多租户标记的字段名及值
}

type _Tenant struct {
	filedName string
	value     any
}

func (t _Tenant) GetTenant() (filedName string, value any) {
	return t.filedName, t.value
}
func (t _Tenant) Data() (data interface{}, err error) {
	out := map[string]any{}
	filed, value := t.GetTenant()
	out[filed] = value
	return out, nil
}

func (t _Tenant) Where() (expressions []goqu.Expression, err error) {
	expressions = make([]goqu.Expression, 0)
	filed, value := t.GetTenant()
	expressions = append(expressions, goqu.C(filed).Eq(value))
	return expressions, nil
}

func NewTenant(filedName string, value any) (tenant _Tenant) {
	return _Tenant{
		filedName: filedName,
		value:     value,
	}
}

func instanceTenantI(tenantI TenantI) _Tenant {
	filed, value := tenantI.GetTenant()
	_tanat := NewTenant(filed, value)
	return _tanat
}

func Insert(tenant TenantI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(instanceTenantI(tenant))
}

func Update(tenant TenantI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(instanceTenantI(tenant))
}

var SorftDelete = Update

func First(tenant TenantI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(instanceTenantI(tenant))
}

func List(tenant TenantI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(instanceTenantI(tenant))
}

func Total(tenant TenantI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(instanceTenantI(tenant))
}

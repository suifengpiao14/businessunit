package tenant

import (
	"strings"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/sqlbuilder"
)

type TenantField sqlbuilder.Field

func (f TenantField) GetTenantField() TenantField {
	return f
}
func (f TenantField) IsSame(o TenantField) bool {
	fv, err := f.ValueFn(nil)
	if err != nil || sqlbuilder.IsNil(fv) {
		return false
	}
	ov, err := f.ValueFn(nil)
	if err != nil || sqlbuilder.IsNil(ov) {
		return false
	}
	return strings.EqualFold(cast.ToString(fv), cast.ToString(ov))
}

type TenantI interface {
	GetTenantField() TenantField // 获取多租户标记的字段名及值
}

func _DataFn(tenantI TenantI) sqlbuilder.DataFn {
	field := tenantI.GetTenantField()
	return sqlbuilder.Field(field).Data
}

func WhereFn(uniqueI TenantI) sqlbuilder.WhereFn {
	return sqlbuilder.Field(uniqueI.GetTenantField()).Where
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

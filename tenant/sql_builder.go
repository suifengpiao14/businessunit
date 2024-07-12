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

func NewTenantField(valueFn sqlbuilder.ValueFn) TenantField {
	field := TenantField{
		Field: *sqlbuilder.NewField(valueFn).SetName("ternat_id").SetTitle("租户ID").MergeSchema(sqlbuilder.Schema{
			Required:  true,
			MinLength: 1,
			MaxLength: 64,
			Maximum:   sqlbuilder.UnsinedInt_maximum_bigint,
			Minimum:   1,
		}),
	}
	return field
}

func (f TenantField) SetName(name string) TenantField {
	f.Field.SetName(name)
	return f
}

func (f TenantField) SetTitle(title string) TenantField {
	f.Field.SetTitle(title)
	return f
}

func (f TenantField) AppendWhereFn(fns ...sqlbuilder.WhereValueFn) TenantField {
	f.Field.WhereFns.Append(fns...)
	return f
}

func (f TenantField) AppendValueFn(fns ...sqlbuilder.ValueFn) TenantField {
	f.Field.ValueFns.Append(fns...)
	return f
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

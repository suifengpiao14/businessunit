package ownerid

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type OwnerIdField struct {
	sqlbuilder.Field
}

func (f OwnerIdField) GetOwnerIdField() OwnerIdField {
	return f
}

type OwnerIdI interface {
	GetOwnerIdField() OwnerIdField
}

var OwnerIdFieldSchema = sqlbuilder.Schema{
	Title:     "所有者",
	Required:  true,
	Comment:   "对象标识,缺失时记录无意义",
	Type:      sqlbuilder.Schema_Type_string,
	MaxLength: 64,
	MinLength: 1,
	Minimum:   1,
}

func NewOwnerIdField(fieldName string, valueFns sqlbuilder.ValueFns, WhereFns sqlbuilder.WhereValueFns, dbSchema *sqlbuilder.Schema) OwnerIdField {
	if dbSchema == nil {
		dbSchema = &OwnerIdFieldSchema
	}
	field := OwnerIdField{
		Field: sqlbuilder.Field{
			Name:   fieldName,
			Schema: dbSchema,
		},
	}
	field.ValueFns.Append(valueFns...)
	field.WhereFns.Append(WhereFns...)
	return field
}

func _DataFn(identityI OwnerIdI) sqlbuilder.DataFn {
	return func() (any, error) {
		field := identityI.GetOwnerIdField()
		if field.ValueFns == nil {
			return nil, nil
		}
		val, err := field.GetValue(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, err
		}

		m := map[string]any{
			sqlbuilder.FieldName2DBColumnName(field.Name): val,
		}
		return m, nil
	}
}

func WhereFn(ownerIdI OwnerIdI) sqlbuilder.WhereFn {

	return ownerIdI.GetOwnerIdField().Where
}

func Insert(ownerIdI OwnerIdI) sqlbuilder.InsertParam {
	// 所有者新增必须写入数据
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(ownerIdI))
}

func Update(ownerIdI OwnerIdI) sqlbuilder.UpdateParam {
	// 所有者不可修改
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(WhereFn(ownerIdI))
}

func First(ownerIdI OwnerIdI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(WhereFn(ownerIdI))
}

func List(ownerIdI OwnerIdI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(WhereFn(ownerIdI))
}

func Total(ownerIdI OwnerIdI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(WhereFn(ownerIdI))
}

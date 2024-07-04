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

func _DataFn(identityI OwnerIdI) sqlbuilder.DataFn {
	return func() (any, error) {
		field := identityI.GetOwnerIdField()
		if field.ValueFn == nil {
			return nil, nil
		}
		val, err := field.ValueFn(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, err
		}

		m := map[string]any{
			field.Name: val,
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

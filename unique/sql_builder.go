package unique

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/softdeleted"
	"github.com/suifengpiao14/sqlbuilder"
)

type UniqueField sqlbuilder.Fields

func (f UniqueField) GetUniqueFields() UniqueField {
	return f
}

type UniqueI interface {
	GetUniqueFields() UniqueField
	sqlbuilder.TableI
	AlreadyExists(sql string) (exists bool, err error)
}

type UniqueIForUpdate interface {
	UniqueI
	identity.IdentityI
}

func _whereFn(uniqueI UniqueI) sqlbuilder.WhereFn {
	return sqlbuilder.Fields(uniqueI.GetUniqueFields()).Where
}

func _checkExists(uniqueI UniqueI, wheres ...sqlbuilder.WhereI) sqlbuilder.DataFn {
	return func() (any, error) {
		totalParam := sqlbuilder.NewTotalBuilder(uniqueI).AppendWhere(wheres...)
		if softdeletedI, ok := uniqueI.(softdeleted.SoftDeletedI); ok { // 如果实现了软删除接口，则排除软删除记录
			totalParam = totalParam.Merge(softdeleted.Total(softdeletedI))
		}
		expressions, err := totalParam.Where()
		if err != nil {
			return nil, err
		}
		if expressions.IsEmpty() { // 条件为空，说明不需要查询唯一索引情况(如产品约定唯一索引字段不能更新)
			return nil, nil
		}

		sql, err := totalParam.ToSQL()
		if err != nil {
			return nil, err
		}
		exists, err := uniqueI.AlreadyExists(sql)
		if err != nil {
			return nil, err
		}
		if exists {
			err = errors.Errorf("unique exists:%s", sqlbuilder.Fields(uniqueI.GetUniqueFields()).String())
			return nil, err
		}
		return nil, err
	}

}

func Insert(uniqueI UniqueI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_checkExists(uniqueI))
}

func Update(uniqueIForUpdate UniqueIForUpdate) sqlbuilder.UpdateParam {
	// 增加排除当前记录
	whereNotID := sqlbuilder.WhereFn(func() (expressions sqlbuilder.Expressions, err error) {
		identity := uniqueIForUpdate.GetIdentityField()
		val, err := identity.GetValue(nil)
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryParseExpressions(sqlbuilder.FieldName2DBColumnName(identity.Name), val); ok {
			return ex, nil
		}
		return sqlbuilder.ConcatExpression(goqu.C(sqlbuilder.FieldName2DBColumnName(identity.Name)).Neq(val)), nil
	})
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_checkExists(uniqueIForUpdate, whereNotID))
}

func First(uniqueI UniqueI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_whereFn(uniqueI))
}

func List(uniqueI UniqueI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_whereFn(uniqueI))
}

func Total(uniqueI UniqueI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_whereFn(uniqueI))
}

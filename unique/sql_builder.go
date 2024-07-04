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
	sqlbuilder.Table
	AlreadyExists(sql string) (exists bool, err error)
}

type UniqueIForUpdate interface {
	UniqueI
	identity.IdentityI
}

func _whereFn(uniqueI UniqueI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		fields := uniqueI.GetUniqueFields()
		expressions = make([]goqu.Expression, 0)
		for _, field := range fields {
			if field.WhereValueFn == nil {
				continue
			}
			val, err := field.WhereValueFn(nil)
			if err != nil {
				return nil, err
			}
			expressions = append(expressions, goqu.C(field.Name).Eq(val))
		}
		return expressions, nil
	}
}

func _checkExists(uniqueI UniqueI, wheres ...sqlbuilder.WhereI) sqlbuilder.DataFn {
	return func() (any, error) {
		totalInstance := _TotalInstance{
			UniqueI: uniqueI,
		}
		totalParam := sqlbuilder.NewTotalBuilder(totalInstance).AppendWhere(wheres...)
		if softdeletedI, ok := uniqueI.(softdeleted.SoftDeletedI); ok { // 如果实现了软删除接口，则排除软删除记录
			totalParam = totalParam.Merge(softdeleted.Total(softdeletedI))
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

type _TotalInstance struct {
	UniqueI
}

func (ins _TotalInstance) Where() (expressions []goqu.Expression, err error) {
	return _whereFn(ins).Where()
}

func Insert(uniqueI UniqueI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_checkExists(uniqueI))
}

func Update(uniqueIForUpdate UniqueIForUpdate) sqlbuilder.UpdateParam {
	// 增加排除当前记录
	whereNotID := sqlbuilder.WhereFn(func() (expressions []goqu.Expression, err error) {
		identity := uniqueIForUpdate.GetIdentityField()
		val, err := identity.ValueFn(nil)
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryParseExpressions(identity.Name, val); ok {
			return ex, nil
		}
		return sqlbuilder.ConcatExpression(goqu.C(identity.Name).Neq(val)), nil
	})
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_checkExists(uniqueIForUpdate, whereNotID))
}

func First(uniqueI UniqueI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(uniqueI UniqueI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(uniqueI UniqueI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

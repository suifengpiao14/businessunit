package softdeleted

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format
var Is_where_EQ = true

type SoftDeletedField struct {
	SoftDeleted sqlbuilder.Field
}

type SoftDeletedI interface {
	GetOperatorField() SoftDeletedField
}

func _SoftDeletedFn(softDeletedI SoftDeletedI) sqlbuilder.DataFn {
	return func() (any, error) {
		col := softDeletedI.GetOperatorField()
		m := map[string]any{}
		val, err := col.SoftDeleted.Value(time.Now().Local().Format(Time_format))
		if err != nil {
			return nil, err
		}
		m[col.SoftDeleted.Name] = val
		return m, nil
	}
}

func _SoftDeletedWhereFn(softDeletedI SoftDeletedI) sqlbuilder.WhereFn {
	col := softDeletedI.GetOperatorField()
	return func() (expressions []goqu.Expression, err error) {
		var expression goqu.Expression
		if Is_where_EQ {
			val, err := col.SoftDeleted.Value("")
			if err != nil {
				return nil, err
			}
			expression = goqu.C(col.SoftDeleted.Name).Eq(val) // 确保删除字段为空
		} else {
			val, err := col.SoftDeleted.Value(nil)
			if err != nil {
				return nil, err
			}
			expression = goqu.C(col.SoftDeleted.Name).Neq(val) // 确保指定字段不等于 特定值
		}
		return sqlbuilder.ConcatExpression(expression), nil
	}
}

func Insert(softDeletedI SoftDeletedI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil)
}

func Update(softDeletedI SoftDeletedI) sqlbuilder.UpdateParam { // 删除
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_SoftDeletedFn(softDeletedI))
}

func First(softDeletedI SoftDeletedI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_SoftDeletedWhereFn(softDeletedI))
}

func List(softDeletedI SoftDeletedI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_SoftDeletedWhereFn(softDeletedI))
}

func Total(softDeletedI SoftDeletedI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_SoftDeletedWhereFn(softDeletedI))
}

package softdeleted

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format
var Is_where_EQ = true

type SoftDeletedField struct {
	SoftDeleted sqlbuilder.Column
}

type SoftDeletedI interface {
	GetOperatorField() SoftDeletedField
}

func _SoftDeletedFn(softDeletedI SoftDeletedI) sqlbuilder.DataFn {
	col := softDeletedI.GetOperatorField()
	m := map[string]any{}
	m[col.SoftDeleted.Name] = col.SoftDeleted.Value(time.Now().Local().Format(Time_format))
	return func() (any, error) {
		return m, nil
	}
}

func _SoftDeletedWhereFn(softDeletedI SoftDeletedI) sqlbuilder.WhereFn {
	col := softDeletedI.GetOperatorField()
	return func() (expressions []goqu.Expression, err error) {
		var expression goqu.Expression
		if Is_where_EQ {
			expression = goqu.C(col.SoftDeleted.Name).Eq(col.SoftDeleted.Value("")) // 确保删除字段为空
		} else {
			expression = goqu.C(col.SoftDeleted.Name).Neq(col.SoftDeleted.Value(nil)) // 确保指定字段不等于 特定值
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

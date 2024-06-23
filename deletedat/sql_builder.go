package deletedat

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type DeletedAtField struct {
	DeletedAt sqlbuilder.Column
}

type DeletedAtI interface {
	GetOperatorField() DeletedAtField
}

func _DeletedAtFn(deletedAtI DeletedAtI) sqlbuilder.DataFn {
	col := deletedAtI.GetOperatorField()
	m := map[string]any{}
	m[col.DeletedAt.Name] = col.DeletedAt.Value(time.Now().Local().Format(Time_format))
	return func() (any, error) {
		return m, nil
	}
}

func _DeletedAtWhereFn(deletedAtI DeletedAtI) sqlbuilder.WhereFn {
	col := deletedAtI.GetOperatorField()
	return func() (expressions []goqu.Expression, err error) {
		sqlbuilder.ConcatExpression(goqu.C(col.DeletedAt.Name).Eq(col.DeletedAt.Value(""))) // 确保删除字段为空
		return
	}
}

func Insert(deletedAtI DeletedAtI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil)
}

func Update(deletedAtI DeletedAtI) sqlbuilder.UpdateParam { // 删除
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DeletedAtFn(deletedAtI))
}

func First(deletedAtI DeletedAtI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_DeletedAtWhereFn(deletedAtI))
}

func List(deletedAtI DeletedAtI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_DeletedAtWhereFn(deletedAtI))
}

func Total(deletedAtI DeletedAtI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_DeletedAtWhereFn(deletedAtI))
}

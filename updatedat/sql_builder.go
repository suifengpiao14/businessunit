package updatedat

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type UpdatedatField sqlbuilder.Field

func (f UpdatedatField) GetUpdatedatField() UpdatedatField {
	return f
}

type UpdatedatI interface {
	GetUpdatedatField() UpdatedatField // 使用每个包下重命名的类型，具有区分度
}

func _DataFn(updatedatI UpdatedatI) sqlbuilder.DataFn {
	col := updatedatI.GetUpdatedatField()
	tim := time.Now().Local().Format(Time_format)
	return func() (any, error) {
		m := map[string]any{}
		val, err := col.Value(tim)
		if err != nil {
			return nil, err
		}
		m[col.Name] = val
		return m, nil
	}
}

func _WhereFn(updatedatI UpdatedatI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := updatedatI.GetUpdatedatField()
		expressions = make([]goqu.Expression, 0)
		val, err := field.Value(nil)
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryConvert2Expressions(val); ok {
			return ex, nil
		}
		if ex, ok := sqlbuilder.TryConvert2Betwwen(field.Name, val); ok {
			return ex, nil
		}
		return expressions, nil
	}
}

func Insert(updatedatI UpdatedatI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(updatedatI))
}

func Update(updatedatI UpdatedatI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(updatedatI))
}

func First(updatedatI UpdatedatI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(updatedatI))
}

func List(updatedatI UpdatedatI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(updatedatI))
}

func Total(updatedatI UpdatedatI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(updatedatI))
}

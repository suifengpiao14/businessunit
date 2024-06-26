package createdat

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type CreatedAtField sqlbuilder.Field

func (f CreatedAtField) GetCreatedAtField() CreatedAtField {
	return f
}

type CreatedAtI interface {
	GetCreatedAtField() CreatedAtField
}

func _DataFn(createdAtI CreatedAtI) sqlbuilder.DataFn {
	return func() (any, error) {
		tim := time.Now().Local().Format(Time_format)
		col := createdAtI.GetCreatedAtField()
		val, err := col.Value(tim)
		if err != nil {
			return nil, err
		}
		m := map[string]any{
			col.Name: val,
		}
		return m, nil
	}
}

func _WhereFn(createdAtI CreatedAtI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := createdAtI.GetCreatedAtField()
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

func Insert(createdAtI CreatedAtI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(createdAtI))
}

func Update(createdAtI CreatedAtI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil)
}

func First(createdAtI CreatedAtI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(createdAtI))
}

func List(createdAtI CreatedAtI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(createdAtI))
}

func Total(createdAtI CreatedAtI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(createdAtI))
}

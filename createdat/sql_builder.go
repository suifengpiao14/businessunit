package createdat

import (
	"time"

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
		field := createdAtI.GetCreatedAtField()
		if field.ValueFn == nil {
			return nil, nil
		}
		val, err := field.ValueFn(tim)
		if err != nil {
			return nil, err
		}
		m := map[string]any{
			field.Name: val,
		}
		return m, nil
	}
}

func _WhereFn(createdAtI CreatedAtI) sqlbuilder.WhereFn {
	return func() (expressions sqlbuilder.Expressions, err error) {
		field := createdAtI.GetCreatedAtField()
		expressions = make(sqlbuilder.Expressions, 0)
		if field.WhereValueFn == nil {
			return nil, err
		}
		val, err := field.WhereValueFn(nil)
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryParseExpressions(field.Name, val); ok {
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

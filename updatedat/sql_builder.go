package updatedat

import (
	"time"

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
		if col.ValueFns == nil {
			return nil, nil
		}
		m := map[string]any{}
		val, err := sqlbuilder.Field(col).GetValue(tim)
		if err != nil {
			return nil, err
		}
		m[sqlbuilder.FieldName2DBColumnName(col.Name)] = val
		return m, nil
	}
}

func _WhereFn(updatedatI UpdatedatI) sqlbuilder.WhereFn {
	return sqlbuilder.Field(updatedatI.GetUpdatedatField()).Where
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

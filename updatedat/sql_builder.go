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
	GetUpdatedatField() UpdatedatField
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

func Insert(updatedatI UpdatedatI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(updatedatI))
}

func Update(updatedatI UpdatedatI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(updatedatI))
}

func First(updatedatI UpdatedatI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(updatedatI UpdatedatI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(updatedatI UpdatedatI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

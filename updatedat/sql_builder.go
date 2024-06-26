package updatedat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type UpdatedatField struct {
	UpdatedAt sqlbuilder.Field
}

type UpdatedatI interface {
	GetUpdatedatField() UpdatedatField
}

func _UpdatedAt(updatedatI UpdatedatI) sqlbuilder.DataFn {
	col := updatedatI.GetUpdatedatField()
	tim := time.Now().Local().Format(Time_format)
	return func() (any, error) {
		m := map[string]any{}
		val, err := col.UpdatedAt.Value(tim)
		if err != nil {
			return nil, err
		}
		m[col.UpdatedAt.Name] = val
		return m, nil
	}
}

func Insert(updatedatI UpdatedatI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_UpdatedAt(updatedatI))
}

func Update(updatedatI UpdatedatI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_UpdatedAt(updatedatI))
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

package createdat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type CreatedAtField struct {
	CreatedAt sqlbuilder.Field
}

type CreatedAtI interface {
	GetCreatedAtField() CreatedAtField
}

func _CreatedAt(createdAtI CreatedAtI) sqlbuilder.DataFn {
	return func() (any, error) {
		tim := time.Now().Local().Format(Time_format)
		col := createdAtI.GetCreatedAtField()
		val, err := col.CreatedAt.Value(tim)
		if err != nil {
			return nil, err
		}
		m := map[string]any{
			col.CreatedAt.Name: val,
		}
		return m, nil
	}
}

func Insert(createdAtI CreatedAtI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_CreatedAt(createdAtI))
}

func Update(createdAtI CreatedAtI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil)
}

func First(createdAtI CreatedAtI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(createdAtI CreatedAtI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(createdAtI CreatedAtI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

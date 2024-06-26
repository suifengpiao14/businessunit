package createdat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type AutoTimeField struct {
	CreatedAt sqlbuilder.Field
}

type AutoTimeI interface {
	GetAutoTimeField() AutoTimeField
}

func _CreatedAt(autoTimeI AutoTimeI) sqlbuilder.DataFn {
	return func() (any, error) {
		tim := time.Now().Local().Format(Time_format)
		col := autoTimeI.GetAutoTimeField()
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

func Insert(autoTime AutoTimeI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_CreatedAt(autoTime))
}

func Update(autoTime AutoTimeI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil)
}

func First(autoTime AutoTimeI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(autoTime AutoTimeI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(autoTime AutoTimeI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

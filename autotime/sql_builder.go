package autotime

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

const (
	Time_format = "2024-01-02 15:04:05"
)

type AutoTimeField struct {
	CreatedAt sqlbuilder.Column
	UpdatedAt sqlbuilder.Column
}

type AutoTimeI interface {
	GetAutoTimeField() AutoTimeField
}

func _CreatedAt(autoTimeI AutoTimeI) sqlbuilder.DataFn {
	tim := time.Now().Local().Format(Time_format)
	col := autoTimeI.GetAutoTimeField()
	m := map[string]any{
		col.CreatedAt.Name: col.CreatedAt.Value(tim),
	}
	return func() (any, error) {
		return m, nil
	}
}

func _UpdatedAt(autoTimeI AutoTimeI) sqlbuilder.DataFn {
	col := autoTimeI.GetAutoTimeField()
	if col.UpdatedAt.Name == "" { // 更新字段可选
		return func() (any, error) {
			return nil, nil
		}
	}
	tim := time.Now().Local().Format(Time_format)

	m := map[string]any{}
	m[col.UpdatedAt.Name] = col.UpdatedAt.Value(tim)
	return func() (any, error) {
		return m, nil
	}
}

func Insert(autoTime AutoTimeI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_CreatedAt(autoTime), _UpdatedAt(autoTime))
}

func Update(autoTime AutoTimeI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_UpdatedAt(autoTime))
}

var SorftDelete = Update

func First(autoTime AutoTimeI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(autoTime AutoTimeI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(autoTime AutoTimeI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

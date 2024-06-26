package autotime

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type AutoTimeField struct {
	CreatedAt sqlbuilder.Field
	UpdatedAt sqlbuilder.Field
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

func _UpdatedAt(autoTimeI AutoTimeI) sqlbuilder.DataFn {
	col := autoTimeI.GetAutoTimeField()
	if col.UpdatedAt.Name == "" { // 更新字段可选
		return func() (any, error) {
			return nil, nil
		}
	}
	tim := time.Now().Local().Format(Time_format)
	return func() (any, error) {
		m := map[string]any{}
		val, err := col.CreatedAt.Value(tim)
		if err != nil {
			return nil, err
		}
		m[col.UpdatedAt.Name] = val
		return m, nil
	}
}

func Insert(autoTime AutoTimeI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_CreatedAt(autoTime), _UpdatedAt(autoTime))
}

func Update(autoTime AutoTimeI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_UpdatedAt(autoTime))
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

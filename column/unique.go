package column

import (
	"github.com/pkg/errors"

	"github.com/suifengpiao14/sqlbuilder"
)

type UniqueColumnI interface {
	AlreadyExists(countSql string) (exists bool, err error)
	sqlbuilder.Table
}

type UniqueColumn struct {
	sqlbuilder.ColumnIs
	UniqueColumnI
}

func (uc UniqueColumn) InsertValue() (value any, err error) {
	val, err := uc.InsertValue()
	if err != nil {
		return nil, err
	}
	_ = val
	return
}

func (uc UniqueColumn) UpdateWhere() (expressions sqlbuilder.Expressions, err error) {
	expressions = make(sqlbuilder.Expressions, 0)
	totalSql, err := sqlbuilder.Total(uc)
	if err != nil {
		return nil, err
	}
	if uc.UniqueColumnI == nil {
		err = errors.Errorf("UniqueColumn.UniqueColumnI required")
		return nil, err
	}
	exists, err := uc.AlreadyExists(totalSql)
	if err != nil {
		return nil, err
	}
	if exists {
		updateWhere, _ := uc.UpdateWhere()
		err = errors.Errorf("unique exists:%s", sqlbuilder.Expression2String(updateWhere...))
		return nil, err
	}
	return expressions, nil
}

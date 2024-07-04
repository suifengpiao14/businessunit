package column

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type UpdatedAtColumn sqlbuilder.Column

func (c UpdatedAtColumn) InsertValue() (value any, err error) {
	if c.InsertValueFn != nil {
		return c.InsertValueFn()
	}
	tim := time.Now().Local().Format(Time_format)
	return tim, nil
}

func (c UpdatedAtColumn) UpdateValue() (value any, err error) {
	if c.UpdateValueFn != nil {
		return c.UpdateValueFn()
	}
	tim := time.Now().Local().Format(Time_format)
	return tim, nil
}

package updatedat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

func NewUpdatedAtField() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(Time_format), nil
	})
	f.SetName("updated_at").SetTitle("更新时间")
	return f
}

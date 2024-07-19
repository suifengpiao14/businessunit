package updatedat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

func OptionUpdatedAt(f *sqlbuilder.Field) {
	f.SetName("updated_at").SetTitle("更新时间")
	f.ValueFns.InsertAsFirst(func(in any) (any, error) {
		return time.Now().Local().Format(Time_format), nil
	})
}

package createdat

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

func OptionCreatedAt(f *sqlbuilder.Field) {
	f.SetName("created_at").SetTitle("创建时间")
	f.SceneInsert(func(f *sqlbuilder.Field) {
		f.ValueFns.InsertAsFirst(func(in any) (any, error) {
			return time.Now().Local().Format(Time_format), nil
		})
	})
	f.SceneUpdate(func(f *sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield) // 更新时屏蔽
	})
}

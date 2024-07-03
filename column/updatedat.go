package column

import (
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

var UpdatedAtColumn = sqlbuilder.Column{
	Name: "updated_at",
	InsertData: func() (value any, err error) {
		tim := time.Now().Local().Format(Time_format)
		return tim, nil
	},
	UpdateData: func() (value any, err error) {
		tim := time.Now().Local().Format(Time_format)
		return tim, nil
	},
}

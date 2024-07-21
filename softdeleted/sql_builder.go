package softdeleted

import (
	"strings"
	"time"

	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

func NewSoftDeletedField(valueType ValueType) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(Time_format), nil
	})
	f.SceneInsert(func(f *sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield)
	})
	f.SceneUpdate(func(f *sqlbuilder.Field) {
		f.ShieldUpdate(true)
		f.WhereFns.Append(func(in any) (any, error) {
			if valueType == ValueType_Delete {
				return sqlbuilder.Neq(in), nil
			}
			return in, nil
		})
	})
	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.WhereFns.Append(func(in any) (any, error) {
			if valueType == ValueType_Delete {
				return sqlbuilder.Neq(in), nil
			}
			return in, nil
		})
	})

	return f
}

type ValueType string

func (vt ValueType) Is(target ValueType) (ok bool) {
	ok = strings.EqualFold(string(vt), string(target))
	return ok
}

const (
	ValueType_Delete ValueType = "delete" //标记 GetDeletedAtField().Value("") 返回的值表示删除值 比如 status=2 表示删除的记录
	ValueType_OK     ValueType = "ok"     //标记 GetDeletedAtField().Value("") 返回的值表示正常记录的值,比如 delted_at="" 表示正常的记录
)

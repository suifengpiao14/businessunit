package enum

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/sqlbuilder"
)

type EnumField struct {
	sqlbuilder.Field
	EnumTitles EnumTitles
}

func (f EnumField) GetEnumField() EnumField {
	return f
}

type EnumTitle struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type EnumTitles []EnumTitle

func (ets EnumTitles) String() string {
	arr := make([]string, 0)
	for _, et := range ets {
		str := fmt.Sprintf("%v-%s", et.Key, et.Title)
		arr = append(arr, str)
	}
	out := strings.Join(arr, ",")
	return out
}

func (enumTitle EnumTitle) IsSame(key string) bool {
	return strings.EqualFold(key, enumTitle.Key)
}

type EnumI interface {
	GetEnumField() EnumField
}

func _DataFn(enumI EnumI) sqlbuilder.DataFn {
	return func() (any, error) {
		col := enumI.GetEnumField()
		if col.ValueFns == nil {
			return nil, nil
		}
		val, err := col.GetValue(nil)
		if err != nil {
			return nil, err
		}
		valid := false
		for _, enum := range col.EnumTitles {
			if strings.EqualFold(cast.ToString(val), enum.Key) {
				valid = true
				break
			}
		}
		if !valid {
			err = errors.Errorf("invalid value except:%s,got:%v", col.EnumTitles.String(), val)
			return nil, err
		}
		m := map[string]any{
			sqlbuilder.FieldName2DBColumnName(col.Name): val,
		}
		return m, nil
	}
}

func WhereFn(enumI EnumI) sqlbuilder.WhereFn {
	field := enumI.GetEnumField()
	return sqlbuilder.Field(field.Field).Where
}

func Insert(enumI EnumI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(enumI))
}

func Update(enumI EnumI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(enumI)).AppendWhere(WhereFn(enumI))
}

func First(enumI EnumI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(WhereFn(enumI))
}

func List(enumI EnumI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(WhereFn(enumI))
}

func Total(enumI EnumI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(WhereFn(enumI))
}

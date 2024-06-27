package enum

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

type EnumField sqlbuilder.Field

func (f EnumField) GetEnumField() EnumField {
	return f
}

type EnumTitle struct {
	Const any    `json:"const"`
	Title string `json:"title"`
}

type EnumTitles []EnumTitle

func (ets EnumTitles) String() string {
	arr := make([]string, 0)
	for _, et := range ets {
		str := fmt.Sprintf("%v-%s", et.Const, et.Title)
		arr = append(arr, str)
	}
	out := strings.Join(arr, ",")
	return out
}

type EnumI interface {
	GetEnumField() EnumField
	EnumTitles() EnumTitles
	IsEqConst(firstConst any, secondConst any) bool // 比较2个常量是否一致
}

func _DataFn(enumI EnumI) sqlbuilder.DataFn {
	return func() (any, error) {
		col := enumI.GetEnumField()
		val, err := col.Value(nil)
		if err != nil {
			return nil, err
		}
		enums := enumI.EnumTitles()
		valid := false
		for _, enum := range enums {
			if enumI.IsEqConst(val, enum.Const) {
				valid = true
				break
			}
		}
		if !valid {
			err = errors.Errorf("invalid value except:%s,got:%v", enums.String(), val)
			return nil, err
		}

		return nil, nil
	}
}

func Insert(enumI EnumI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(enumI))
}

func Update(enumI EnumI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(enumI))
}

func First(enumI EnumI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(enumI EnumI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(enumI EnumI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}

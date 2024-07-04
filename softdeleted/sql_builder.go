package softdeleted

import (
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

type SoftDeletedField sqlbuilder.Field

func (f SoftDeletedField) GetDeletedAtField() SoftDeletedField {
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

type SoftDeletedI interface {
	GetDeletedAtField() (valueType ValueType, softDeletedField SoftDeletedField)
}

func _DataFn(softDeletedI SoftDeletedI) sqlbuilder.DataFn {
	return func() (any, error) {
		_, field := softDeletedI.GetDeletedAtField()
		m := map[string]any{}
		val, err := field.ValueFn(time.Now().Local().Format(Time_format))
		if err != nil {
			return nil, err
		}
		m[field.Name] = val
		return m, nil
	}
}

func _WhereFn(softDeletedI SoftDeletedI) sqlbuilder.WhereFn {
	valueType, field := softDeletedI.GetDeletedAtField()
	return func() (expressions []goqu.Expression, err error) {
		if field.ValueFn == nil {
			return nil, nil
		}
		val, err := field.WhereValueFn("")
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryConvert2Expressions(val); ok {
			return ex, nil
		}
		var expression goqu.Expression
		switch valueType {
		case ValueType_OK:
			expression = goqu.C(field.Name).Eq(val) // 确保删除字段为空
		case ValueType_Delete:
			expression = goqu.C(field.Name).Neq(val) // 确保指定字段不等于 特定值
		default:
			err = errors.Errorf("invalid valueType , except %s|%s,got:%s", ValueType_OK, ValueType_Delete, valueType)
			return nil, err
		}
		return sqlbuilder.ConcatExpression(expression), nil
	}
}

func Insert(softDeletedI SoftDeletedI) sqlbuilder.InsertParam { // softdelete 没有insert场景，此处仅仅补齐，方便集成
	return sqlbuilder.NewInsertBuilder(nil)
}

func Delete(softDeletedI SoftDeletedI) sqlbuilder.UpdateParam { // 删除
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(softDeletedI))
}
func Update(softDeletedI SoftDeletedI) sqlbuilder.UpdateParam { // 这个地方方便，和其它单元集成，改成条件，需要删除数据，使用Delete方法
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(_WhereFn(softDeletedI))
}

func First(softDeletedI SoftDeletedI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(softDeletedI))
}

func List(softDeletedI SoftDeletedI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(softDeletedI))
}

func Total(softDeletedI SoftDeletedI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(softDeletedI))
}

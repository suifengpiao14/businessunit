package unique

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

type CheckUniqueI interface {
	sqlbuilder.TableI
	AlreadyExists(sql string) (exists bool, err error) //
}

func _checkExists(uniqueI CheckUniqueI, uniqueField sqlbuilder.Fields, idField *sqlbuilder.Field) {
	if len(uniqueField) == 0 {
		return
	}
	firstField := uniqueField[0] // 先获取原始的第一个元素
	//  在第一列中，增加查询唯一键是否存在条件,放到valueFns 内部，新增、修改都能检测到,放到whereFn, insert 时检测不到
	firstField.ValueFns.Append(func(in any) (any, error) { // 放到闭包函数内，延迟执行
		uniqueField = uniqueField.Copy() // 复制一份数据，不影响外部
		uniqueField.AppendWhereValueFn(sqlbuilder.ValueFnForward)
		if idField != nil {
			idField = idField.Copy() // 复制一份，不影响外部的
		}
		totalParam := sqlbuilder.NewTotalBuilder(uniqueI.Table()).AppendFields(uniqueField...)
		expressions, err := totalParam.Where()
		if err != nil {
			return nil, err
		}
		if expressions.IsEmpty() { // 条件为空，说明不需要查询唯一索引情况(如产品约定唯一索引字段不能更新)
			return nil, nil
		}
		if idField != nil {
			idField.WhereFns.Append(func(in any) (any, error) {
				return goqu.C(idField.DBName()).Neq(in), nil
			})
			totalParam = totalParam.AppendFields(idField)
		}

		sql, err := totalParam.ToSQL()
		if err != nil {
			return nil, err
		}
		exists, err := uniqueI.AlreadyExists(sql)
		if err != nil {
			return nil, err
		}
		if exists {
			err = errors.Errorf("unique exists:%s", sqlbuilder.Fields(uniqueField).String())
			return nil, err
		}
		return in, nil

	})

}

func OptionUnique(checkUniqueI CheckUniqueI, idField *sqlbuilder.Field) func(fields ...*sqlbuilder.Field) {
	return func(fields ...*sqlbuilder.Field) {
		if len(fields) == 0 {
			return
		}
		first := fields[0]
		first.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			_checkExists(checkUniqueI, fields, nil)
		})
		first.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			_checkExists(checkUniqueI, fields, idField)
		})
	}
}

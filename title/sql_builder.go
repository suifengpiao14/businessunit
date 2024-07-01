package title

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/sqlbuilder"
)

type TitleI interface {
	GetTitleField() sqlbuilder.Field
	identity.IdentityI
}

func _DataFn(titleI TitleI) sqlbuilder.DataFn {
	return func() (any, error) {
		title := titleI.GetTitleField()
		m := map[string]any{}
		val, err := title.Value(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, nil
		}
		m[title.Name] = val
		return m, nil
	}
}
func _WhereFn(titleI TitleI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := titleI.GetIdentityField()
		expressions = make([]goqu.Expression, 0)
		val, err := field.Value(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, nil
		}
		if ex, ok := sqlbuilder.TryConvert2Expressions(val); ok {
			return ex, nil
		}
		likeValue := "%" + cast.ToString(val) + "%"
		expressions = append(expressions, goqu.C(field.Name).ILike(likeValue))
		return expressions, nil
	}
}

func Insert(titleI TitleI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).Merge(identity.Insert(titleI)).AppendData(_DataFn(titleI))
}

func Update(titleI TitleI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).Merge(identity.Update(titleI)).AppendData(_DataFn(titleI))
}

func First(titleI TitleI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).Merge(identity.First(titleI)).AppendWhere(_WhereFn(titleI))
}

func List(titleI TitleI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).Merge(identity.List(titleI)).AppendWhere(_WhereFn(titleI))
}

func Total(titleI TitleI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).Merge(identity.Total(titleI)).AppendWhere(_WhereFn(titleI))
}

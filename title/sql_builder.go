package title

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/sqlbuilder"
)

type Title struct {
	Title sqlbuilder.Field
	ID    identity.IdentityField
}

func (t Title) GetIdentityField() identity.IdentityField {
	return t.ID
}

func (t Title) GetTitle() Title {
	return t
}

type TitleI interface {
	GetTitle() Title
}

func _DataFn(titleI TitleI) sqlbuilder.DataFn {
	return func() (any, error) {
		title := titleI.GetTitle()
		if title.Title.ValueFn == nil {
			return nil, nil
		}
		m := map[string]any{}
		val, err := title.Title.ValueFn(nil)
		if err != nil {
			return nil, err
		}
		if sqlbuilder.IsNil(val) {
			return nil, nil
		}
		m[title.Title.Name] = val
		return m, nil
	}
}
func _WhereFn(titleI TitleI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		field := titleI.GetTitle().ID
		expressions = make([]goqu.Expression, 0)
		if field.WhereValueFn == nil {
			return nil, nil
		}
		val, err := field.WhereValueFn(nil)
		if err != nil {
			return nil, err
		}
		if ex, ok := sqlbuilder.TryParseExpressions(field.Name, val); ok {
			return ex, nil
		}
		likeValue := "%" + cast.ToString(val) + "%"
		expressions = append(expressions, goqu.C(field.Name).ILike(likeValue))
		return expressions, nil
	}
}

func Insert(titleI TitleI) sqlbuilder.InsertParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewInsertBuilder(nil).Merge(identity.Insert(title)).AppendData(_DataFn(titleI))
}

func Update(titleI TitleI) sqlbuilder.UpdateParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewUpdateBuilder(nil).Merge(identity.Update(title)).AppendData(_DataFn(titleI))
}

func First(titleI TitleI) sqlbuilder.FirstParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewFirstBuilder(nil).Merge(identity.First(title)).AppendWhere(_WhereFn(titleI))
}

func List(titleI TitleI) sqlbuilder.ListParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewListBuilder(nil).Merge(identity.List(title)).AppendWhere(_WhereFn(titleI))
}

func Total(titleI TitleI) sqlbuilder.TotalParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewTotalBuilder(nil).Merge(identity.Total(title)).AppendWhere(_WhereFn(titleI))
}

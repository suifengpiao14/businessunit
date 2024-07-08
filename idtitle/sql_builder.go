package idtitle

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/sqlbuilder"
)

type IdTitle struct {
	Title sqlbuilder.Field
	ID    identity.IdentityField
}

func (t IdTitle) GetIdentityField() identity.IdentityField {
	return t.ID
}

func (t IdTitle) GetTitle() IdTitle {
	return t
}

var TitleFieldSchema = sqlbuilder.DBSchema{
	Required:  false,
	Type:      sqlbuilder.DBSchema_Type_string,
	MaxLength: 64,
	MinLength: 1,
}

// NewTitleField 生成标题列，标题类一般没有逻辑，主要用于配合ID显示
func NewTitleField(valueFn sqlbuilder.ValueFn) sqlbuilder.Field {
	field := sqlbuilder.NewField(valueFn).SetName("title").MergeDBSchema(TitleFieldSchema).SetTitle("标题")
	return field
}

// NewTitleField 生成标题列，标题类一般没有逻辑，主要用于配合ID显示
var NewIdentityField = identity.NewIdentityField

func NewIdTitleFiled(idValueFn sqlbuilder.ValueFn, titleValueFn sqlbuilder.ValueFn) IdTitle {
	return IdTitle{
		ID:    NewIdentityField(idValueFn),
		Title: NewTitleField(titleValueFn),
	}
}

type IdTitleI interface {
	GetTitle() IdTitle
}

func _DataFn(titleI IdTitleI) sqlbuilder.DataFn {
	return func() (any, error) {
		title := titleI.GetTitle()
		if title.Title.ValueFns == nil {
			return nil, nil
		}
		m := map[string]any{}
		val, err := title.Title.GetValue(nil)
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
func _WhereFn(titleI IdTitleI) sqlbuilder.WhereFn {
	return func() (expressions sqlbuilder.Expressions, err error) {
		field := titleI.GetTitle().Title
		field.WhereFns.Insert(-1, func(in any) (value any, err error) {
			val := cast.ToString(in)
			if val == "" {
				return in, nil
			}
			likeValue := "%" + val + "%"
			expressions = append(expressions, goqu.C(field.Name).ILike(likeValue))
			return expressions, nil
		})

		return field.Where()
	}
}

func Insert(titleI IdTitleI) sqlbuilder.InsertParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewInsertBuilder(nil).Merge(identity.Insert(title)).AppendData(_DataFn(titleI))
}

func Update(titleI IdTitleI) sqlbuilder.UpdateParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewUpdateBuilder(nil).Merge(identity.Update(title)).AppendData(_DataFn(titleI))
}

func First(titleI IdTitleI) sqlbuilder.FirstParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewFirstBuilder(nil).Merge(identity.First(title)).AppendWhere(_WhereFn(titleI))
}

func List(titleI IdTitleI) sqlbuilder.ListParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewListBuilder(nil).Merge(identity.List(title)).AppendWhere(_WhereFn(titleI))
}

func Total(titleI IdTitleI) sqlbuilder.TotalParam {
	title := titleI.GetTitle()
	return sqlbuilder.NewTotalBuilder(nil).Merge(identity.Total(title)).AppendWhere(_WhereFn(titleI))
}

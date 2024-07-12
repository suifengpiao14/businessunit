package idtitle

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/sqlbuilder"
)

type IdTitleField struct {
	ID    *identity.IdentityField
	Title *sqlbuilder.Field
}

func (t IdTitleField) GetIdentityField() *identity.IdentityField {
	return t.ID
}

func (t *IdTitleField) GetIdTitle() *IdTitleField {
	return t
}
func (t IdTitleField) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{*t.Title, t.ID.Field}
}

var TitleFieldSchema = sqlbuilder.Schema{
	Required:  false,
	Type:      sqlbuilder.Schema_Type_string,
	MaxLength: 64,
	MinLength: 1,
}

// NewTitleField 生成标题列，标题类一般没有逻辑，主要用于配合ID显示
func NewTitleField(valueFn sqlbuilder.ValueFn) *sqlbuilder.Field {
	field := sqlbuilder.NewField(valueFn).SetName("title").MergeSchema(TitleFieldSchema).SetTitle("标题")
	return field
}

// NewTitleField 生成标题列，标题类一般没有逻辑，主要用于配合ID显示
var NewIdentityField = identity.NewIdentityField

func NewIdTitleFiled(idValueFn sqlbuilder.ValueFn, titleValueFn sqlbuilder.ValueFn) *IdTitleField {
	return &IdTitleField{
		ID:    NewIdentityField(idValueFn),
		Title: NewTitleField(titleValueFn),
	}
}

type IdTitleI interface {
	GetIdTitle() *IdTitleField
}

func _DataFn(titleI IdTitleI) sqlbuilder.DataFn {
	return func() (any, error) {
		title := titleI.GetIdTitle()
		if title.Title.ValueFns == nil {
			return nil, nil
		}
		m := map[string]any{}
		val, err := title.Title.GetValue(nil)
		if err != nil {
			return nil, err
		}
		m[sqlbuilder.FieldName2DBColumnName(title.Title.Name)] = val
		return m, nil
	}
}
func _WhereFn(titleI IdTitleI) sqlbuilder.WhereFn {
	return func() (expressions sqlbuilder.Expressions, err error) {
		field := titleI.GetIdTitle().Title
		field.WhereFns.Insert(-1, func(dbColumnName string, in any) (value any, err error) {
			val := cast.ToString(in)
			if val == "" {
				return in, nil
			}
			likeValue := "%" + val + "%"
			expressions = append(expressions, goqu.C(sqlbuilder.FieldName2DBColumnName(field.Name)).ILike(likeValue))
			return expressions, nil
		})

		return field.Where()
	}
}

func Insert(titleI IdTitleI) sqlbuilder.InsertParam {
	id := titleI.GetIdTitle().ID
	return sqlbuilder.NewInsertBuilder(nil).Merge(identity.Insert(id)).AppendData(_DataFn(titleI))
}

func Update(titleI IdTitleI) sqlbuilder.UpdateParam {
	id := titleI.GetIdTitle().ID
	return sqlbuilder.NewUpdateBuilder(nil).Merge(identity.Update(id)).AppendData(_DataFn(titleI))
}

func First(titleI IdTitleI) sqlbuilder.FirstParam {
	id := titleI.GetIdTitle().ID
	return sqlbuilder.NewFirstBuilder(nil).Merge(identity.First(id)).AppendWhere(_WhereFn(titleI))
}

func List(titleI IdTitleI) sqlbuilder.ListParam {
	id := titleI.GetIdTitle().ID
	return sqlbuilder.NewListBuilder(nil).Merge(identity.List(id)).AppendWhere(_WhereFn(titleI))
}

func Total(titleI IdTitleI) sqlbuilder.TotalParam {
	id := titleI.GetIdTitle().ID
	return sqlbuilder.NewTotalBuilder(nil).Merge(identity.Total(id)).AppendWhere(_WhereFn(titleI))
}

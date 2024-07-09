package identity

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

type IdentityField struct {
	sqlbuilder.Field
}

func (f *IdentityField) GetIdentityField() *IdentityField {
	return f
}
func (f *IdentityField) SetName(name string) *IdentityField {
	f.Name = name
	return f
}

func (f *IdentityField) SetTitle(title string) *IdentityField {
	f.Field.SetTitle(title)
	return f
}

var IdentityFieldSchema = sqlbuilder.DBSchema{
	Required:  true,
	Type:      sqlbuilder.DBSchema_Type_string,
	MaxLength: 64,
	MinLength: 1,
}

// NewIdentityField 生成标题列，标题类一般没有逻辑，主要用于配合ID显示
func NewIdentityField(valueFn sqlbuilder.ValueFn) *IdentityField {
	field := &IdentityField{
		Field: *sqlbuilder.NewField(valueFn).SetName("id").SetTitle("ID").MergeDBSchema(IdentityFieldSchema),
	}
	field.WhereFns.Append(func(in any) (any, error) {
		return goqu.Ex{
			field.Name: in,
		}, nil
	})
	return field
}

type IdentityI interface {
	GetIdentityField() *IdentityField
}

func _DataFn(identityI IdentityI) sqlbuilder.DataFn {
	return func() (any, error) {
		field := identityI.GetIdentityField()
		if field.ValueFns == nil {
			return nil, nil
		}
		val, err := field.GetValue(nil)
		if err != nil {
			return nil, err
		}

		m := map[string]any{
			field.Name: val,
		}
		return m, nil
	}
}

func _WhereFn(identityI IdentityI) sqlbuilder.WhereFn {
	field := identityI.GetIdentityField()
	return field.Where
}

func Insert(identityI IdentityI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(identityI))
}

func Update(identityI IdentityI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendWhere(_WhereFn(identityI))
}

func First(identityI IdentityI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(identityI))
}

func List(identityI IdentityI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(identityI))
}

func Total(identityI IdentityI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(identityI))
}

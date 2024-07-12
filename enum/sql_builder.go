package enum

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

type EnumField struct {
	sqlbuilder.Field
	EnumTitles sqlbuilder.Enums
}

func (f EnumField) GetEnumField() EnumField {
	return f
}

func NewEnumField(valueFn sqlbuilder.ValueFn, enums ...sqlbuilder.Enum) EnumField {
	schema := sqlbuilder.Schema{
		Enums: enums,
	}
	return EnumField{
		Field: *sqlbuilder.NewField(valueFn).MergeSchema(schema),
	}
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
		valid := col.EnumTitles.Contains(val)
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

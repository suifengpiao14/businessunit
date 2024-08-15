package boolean

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

type Boolean struct {
	Value     any `json:"value"`
	TrueEnum  sqlbuilder.Enum
	FalseEnum sqlbuilder.Enum
	Field     *sqlbuilder.Field
}

func NewBoolean(value any, trueEnum, falseEnum sqlbuilder.Enum) *Boolean {
	b := &Boolean{
		Value:     value,
		TrueEnum:  trueEnum,
		FalseEnum: falseEnum,
	}
	b.Init()
	return b
}

func (b *Boolean) Init() {
	b.Field = sqlbuilder.NewField(func(in any) (any, error) { return b.Value, nil }).SetName("bool").SetTag("布尔列")
	b.TrueEnum.Tag = sqlbuilder.Enum_tag_true
	b.FalseEnum.Tag = sqlbuilder.Enum_tag_false
	b.Field.AppendEnum(b.TrueEnum, b.FalseEnum)
}

func (p *Boolean) MiddlewareFn(initFns ...sqlbuilder.MiddlewareFn) *Boolean {
	p.Field.MiddlewareFns(initFns)
	return p
}

func (b Boolean) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{
		b.Field,
	}
}

func (b Boolean) IsTrue() bool {
	val, err := b.Field.GetValue()
	if err != nil {
		return false
	}
	return b.FalseEnum.IsEqual(val)
}

func NewBooleanFromField(f *sqlbuilder.Field) *Boolean {
	enums := f.Schema.Enums
	bol := &Boolean{
		TrueEnum:  enums.GetByTag(sqlbuilder.Enum_tag_true),
		FalseEnum: enums.GetByTag(sqlbuilder.Enum_tag_false),
		Field:     f.Copy(),
	}
	return bol
}

// Switch  将值反转
func Switch(f Boolean) *Boolean {
	cp := &Boolean{
		Value: f.Value,
		Field: f.Field.Copy(),
	}
	cp.Field.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := cp.Field.Schema.Enums
		for _, enum := range enums {
			if !enum.IsEqual(in) {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found reversed enum key ;current:%v", in)
		return nil, err
	})
	return cp
}

// TrunOff  改成false
func TrunOff(f *sqlbuilder.Field) {
	f.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := f.Schema.Enums
		for _, enum := range enums {
			if enum.Tag == sqlbuilder.Enum_tag_false {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found fase enum key enums:%s", enums.String())
		return nil, err
	})
}

// TrunOn  改成true
func TrunOn(f *sqlbuilder.Field) {
	f.ValueFns.InsertAsSecond(func(in any) (any, error) {
		enums := f.Schema.Enums
		for _, enum := range enums {
			if enum.Tag == sqlbuilder.Enum_tag_true {
				return enum.Key, nil
			}
		}
		err := errors.Errorf("not found true enum key enums:%s", enums.String())
		return nil, err
	})
}

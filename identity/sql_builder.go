package identity

import "github.com/suifengpiao14/sqlbuilder"

type Identifier struct {
	Value string `json:"value"`
	Field *sqlbuilder.Field
}

func (i *Identifier) Init() {
	i.Field = sqlbuilder.NewField(func(in any) (any, error) { return i.Value, nil }).SetName("identity").SetTitle("标识")
	i.Field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
}

func NewIdentifier(val string) *Identifier {
	f := &Identifier{
		Value: val,
	}
	f.Init()
	return f
}

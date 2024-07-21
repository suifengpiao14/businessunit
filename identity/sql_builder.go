package identity

import "github.com/suifengpiao14/sqlbuilder"

func NewIdentityField(valueFn sqlbuilder.ValueFn) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(valueFn).SetName("identity").SetTitle("标识")
	f.SceneSelect(func(f *sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

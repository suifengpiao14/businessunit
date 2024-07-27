package businessunit

import "github.com/suifengpiao14/sqlbuilder"

func NewNickname(nickname string) *sqlbuilder.Field {
	f := NewNameField(nickname).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("nickname").SetTitle("昵称")
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(
			sqlbuilder.ValueFnEmpty2Nil,
			sqlbuilder.ValueFnWhereLike,
		)
	})
	return f
}

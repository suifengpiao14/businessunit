package ownerid

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type OwnerID struct {
	Value any `json:"value"`
	Field *sqlbuilder.Field
}

func NewOwnerID(value any) *OwnerID {
	o := &OwnerID{Value: value}
	o.Init()
	return o
}

func (o *OwnerID) Init() {
	o.Field = sqlbuilder.NewField(func(in any) (any, error) { return o.Value, nil }).SetName("ownerId").SetTitle("所有者").MergeSchema(sqlbuilder.Schema{
		Comment:      "对象标识,缺失时记录无意义",
		Type:         sqlbuilder.Schema_Type_string,
		MaxLength:    64,
		MinLength:    1,
		Minimum:      1,
		ShieldUpdate: true, // 所有者不可跟新
	})
	o.Field.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetRequired(true)
	})
	o.Field.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true)
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	o.Field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.AppendIfNotFirst(sqlbuilder.ValueFnEmpty2Nil)
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
}

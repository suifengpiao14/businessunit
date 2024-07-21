package tenant

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type Tenant struct {
	Value any `json:"value"`
	Field *sqlbuilder.Field
}

func (t *Tenant) Init() {
	t.Field = sqlbuilder.NewField(func(in any) (any, error) { return t.Value, nil }).SetName("ternatId").SetTitle("租户ID")
	t.Field.MergeSchema(sqlbuilder.Schema{
		Required:  true,
		MinLength: 1,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Maximum:   sqlbuilder.UnsinedInt_maximum_bigint,
		Minimum:   1,
	})
	t.Field.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // 不可更新
	})
	t.Field.WhereFns.InsertAsFirst(sqlbuilder.ValueFnForward) // update,select 都必须为条件
}

func NewTenant(value any) *Tenant {
	t := &Tenant{Value: value}
	t.Init()
	return t
}

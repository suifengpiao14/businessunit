package ownerid_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/autoid"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/sqlbuilder"
)

type UpdateParam struct {
	ID      string `db:"-"`
	Name    string
	OwnerID int
}

func (p UpdateParam) Table() string {
	return "t_table"
}

func TestUpdate(t *testing.T) {
	p := UpdateParam{
		ID:      "qwqwgwerst3wyt5y456u56uj5rywr",
		Name:    "张三",
		OwnerID: 14,
	}

	sql, err := sqlbuilder.NewUpdateBuilder(p).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:      "",
		Name:    "张三",
		OwnerID: 1,
	}
	identityField := sqlbuilder.NewField(func(in any) (any, error) { return p.ID, nil }).WithOptions(autoid.Insert)
	ownerIdField := sqlbuilder.NewField(func(in any) (any, error) { return p.OwnerID, nil }).WithOptions(ownerid.Insert)
	nameField := sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).SetName("name")
	sql, err := sqlbuilder.NewInsertBuilder(p).AppendField(identityField, ownerIdField, nameField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

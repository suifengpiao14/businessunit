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

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:      "",
		Name:    "张三",
		OwnerID: 1,
	}
	identityField := sqlbuilder.NewField(func(in any) (any, error) { return p.ID, nil }).WithOptions(autoid.OptionAutoID)
	ownerIdField := sqlbuilder.NewField(func(in any) (any, error) { return p.OwnerID, nil }).WithOptions(ownerid.OptionOwnerID)
	nameField := sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).SetName("name")
	sql, err := sqlbuilder.NewInsertBuilder(p).AppendField(identityField, ownerIdField, nameField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestUpdate(t *testing.T) {
	p := UpdateParam{
		ID:      "",
		Name:    "张三",
		OwnerID: 1,
	}
	identityField := sqlbuilder.NewField(func(in any) (any, error) { return p.ID, nil }).WithOptions(autoid.OptionAutoID)
	ownerIdField := sqlbuilder.NewField(func(in any) (any, error) { return p.OwnerID, nil }).WithOptions(ownerid.OptionOwnerID)
	nameField := sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).SetName("name")

	sql, err := sqlbuilder.NewUpdateBuilder(p).AppendField(identityField, ownerIdField, nameField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

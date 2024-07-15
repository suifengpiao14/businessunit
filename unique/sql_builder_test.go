package unique_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/id"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/businessunit/unique"
	"github.com/suifengpiao14/sqlbuilder"
)

type UpdateParam struct {
	ID   string `db:"-"`
	Name string
}

func (p UpdateParam) Table() string {
	return "t_table"
}

func (p UpdateParam) AlreadyExists(sql string) (exists bool, err error) {
	fmt.Println(sql)
	return false, err
}

func TestUpdate(t *testing.T) {
	p := UpdateParam{
		ID:   "15",
		Name: "张三",
	}

	idField := sqlbuilder.NewField(func(in any) (any, error) { return p.ID, nil }).WithOptions(id.Update)
	uniqueFields := sqlbuilder.NewFields(sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).WithOptions(name.Update)).WithOptions(unique.UpdateFn(p, idField))
	sql, err := sqlbuilder.NewUpdateBuilder(sqlbuilder.TableFn(p.Table)).AppendField(*uniqueFields...).AppendField(idField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:   "15",
		Name: "张三",
	}
	uniqueFields := sqlbuilder.NewFields(sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).WithOptions(name.Update)).WithOptions(unique.Insert(p))
	sql, err := sqlbuilder.NewInsertBuilder(sqlbuilder.TableFn(p.Table)).AppendField(*uniqueFields...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

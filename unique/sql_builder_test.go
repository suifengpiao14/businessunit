package unique_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/autoid"
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

	idField := autoid.NewAutoIdField(func(in any) (any, error) { return p.ID, nil })
	uniqueFields := sqlbuilder.NewFields(sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil })).WithOptions(unique.OptionUnique(p, idField))
	sql, err := sqlbuilder.NewUpdateBuilder(p.Table()).AppendFields(uniqueFields...).AppendFields(idField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:   "15",
		Name: "张三",
	}
	uniqueFields := sqlbuilder.NewFields(sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil })).WithOptions(unique.OptionUnique(p, nil))
	sql, err := sqlbuilder.NewInsertBuilder(p.Table()).AppendFields(uniqueFields...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

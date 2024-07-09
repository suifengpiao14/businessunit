package unique_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/softdeleted"
	"github.com/suifengpiao14/businessunit/unique"
	"github.com/suifengpiao14/sqlbuilder"
)

type UpdateParam struct {
	ID   string `db:"-"`
	Name string
}

func (p UpdateParam) GetIdentityField() *identity.IdentityField {
	return identity.NewIdentityField(func(in any) (any, error) { return p.ID, nil })
}
func (p UpdateParam) GetUniqueFields() (fields unique.UniqueField) {
	fields = make(unique.UniqueField, 0)
	fields = append(fields, sqlbuilder.Field{
		DBName: "Fname",
		ValueFns: sqlbuilder.ValueFns{func(in any) (value any, err error) {
			return p.Name, nil
		}},
	},
	)
	return fields
}

func (p UpdateParam) GetDeletedAtField() (softdeleted.ValueType, softdeleted.SoftDeletedField) {
	return softdeleted.ValueType_OK, softdeleted.SoftDeletedField{
		DBName: "Fdeleted_at",
		ValueFns: sqlbuilder.ValueFns{func(in any) (value any, err error) {
			return "", nil
		}},
	}
}

func (p UpdateParam) Where() (expressions sqlbuilder.Expressions, err error) {
	return nil, nil
}
func (p UpdateParam) Table() string {
	return "t_table"
}
func (p UpdateParam) Data() (data any, err error) {
	return p, nil
}
func (p UpdateParam) AlreadyExists(sql string) (exists bool, err error) {
	fmt.Println(sql)
	return true, err
}

func TestUpdate(t *testing.T) {
	p := UpdateParam{
		ID:   "15",
		Name: "张三",
	}

	sql, err := sqlbuilder.NewUpdateBuilder(p).Merge(unique.Update(p), identity.Update(p), softdeleted.Update(p)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:   "15",
		Name: "张三",
	}

	sql, err := sqlbuilder.NewInsertBuilder(p).Merge(unique.Insert(p), identity.Insert(p), softdeleted.Insert(p)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

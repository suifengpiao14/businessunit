package ownerid_test

import (
	"fmt"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/businessunit/softdeleted"
	"github.com/suifengpiao14/sqlbuilder"
)

type UpdateParam struct {
	ID      string `db:"-"`
	Name    string
	OwnerID int
}

func (p UpdateParam) GetIdentityField() *identity.IdentityField {
	return identity.NewIdentityField(func(in any) (any, error) { return p.ID, nil }).SetName("Fid")
}
func (p UpdateParam) GetOwnerIdField() (field ownerid.OwnerIdField) {
	field = ownerid.NewOwnerIdField(
		"Fonwer_id",
		sqlbuilder.NewValueFns(func() (value any, err error) { return p.OwnerID, nil }),
		nil, nil,
	)
	return field
}

func (p UpdateParam) GetDeletedAtField() (softdeleted.ValueType, softdeleted.SoftDeletedField) {
	return softdeleted.ValueType_OK, softdeleted.SoftDeletedField{
		DBName:   "Fdeleted_at",
		ValueFns: sqlbuilder.NewValueFns(func() (value any, err error) { return "", nil }),
	}
}

func (p UpdateParam) Where() (expressions []goqu.Expression, err error) {
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
		ID:      "qwqwgwerst3wyt5y456u56uj5rywr",
		Name:    "张三",
		OwnerID: 14,
	}

	sql, err := sqlbuilder.NewUpdateBuilder(p).Merge(ownerid.Update(p), identity.Update(p), softdeleted.Update(p)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestInsert(t *testing.T) {
	p := UpdateParam{
		ID:      "15",
		Name:    "张三",
		OwnerID: 1,
	}

	sql, err := sqlbuilder.NewInsertBuilder(p).Merge(ownerid.Insert(p), identity.Insert(p), softdeleted.Insert(p)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

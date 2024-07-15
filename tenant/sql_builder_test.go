package tenant_test

import (
	"fmt"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/businessunit/tenant"
	"github.com/suifengpiao14/businessunit/uuid"
	"github.com/suifengpiao14/sqlbuilder"
)

func init() {
	sqlbuilder.Dialect = sqlbuilder.Dialect_Mysql
}

type InsertParam struct {
	Id     string `db:"Fid"`
	Name   string `db:"Fname"`
	Email  string `db:"Femail"`
	Tenant string `db:"-"`
}

func (p InsertParam) Table() string {
	return "t_user"
}

func (p InsertParam) Fields() sqlbuilder.Fields {
	fields := make(sqlbuilder.Fields, 0)
	fields = append(fields, sqlbuilder.NewField(func(in any) (any, error) { return p.Id, nil }).WithOptions(uuid.Insert))
	fields = append(fields, sqlbuilder.NewField(func(in any) (any, error) { return p.Name, nil }).WithOptions(name.Insert))
	fields = append(fields, sqlbuilder.NewField(func(in any) (any, error) { return p.Email, nil }).WithOptions(email.Insert))
	fields = append(fields, sqlbuilder.NewField(func(in any) (any, error) { return p.Tenant, nil }).WithOptions(tenant.Insert))
	return fields

}

func TestInsert(t *testing.T) {
	row := InsertParam{
		Id:     "123",
		Name:   "张三",
		Email:  "莉丝",
		Tenant: "1000001",
	}
	tenantField := sqlbuilder.NewField(func(in any) (any, error) { return row.Tenant, nil }).WithOptions(tenant.Insert)
	sql, err := sqlbuilder.NewInsertBuilder(row).AppendField(tenantField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

type UpdateParam struct {
	Id    string `db:"-"`
	Name  string `db:"Fname"`
	Email string `db:"Femail"`
}

func (p UpdateParam) Table() string {
	return "t_user"
}

func (p UpdateParam) Validate() (err error) {
	if p.Id == "" {
		err = errors.New("id required")
	}
	return err
}

func (p UpdateParam) Data() (data interface{}, err error) {
	return p, nil
}

func (p UpdateParam) Where() (expressions sqlbuilder.Expressions, err error) {
	expressions = sqlbuilder.Expressions{
		goqu.C("Fid").Eq(p.Id),
		goqu.C("Fdeleted").Eq(""),
	}
	return
}

func TestUpdate(t *testing.T) {
	row := UpdateParam{
		Id:    "123",
		Name:  "张三",
		Email: "莉丝",
	}
	sql, err := sqlbuilder.NewUpdateBuilder(row).Merge(tenant.Update(tenantData)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

type FirstParam struct {
	Id    string `db:"-"`
	Name  string `db:"Fname"`
	Email string `db:"Femail"`
}

func (p FirstParam) Table() string {
	return "t_user"
}

func (p FirstParam) Validate() (err error) {
	if p.Id == "" {
		err = errors.New("id required")
	}
	return err
}

func (p FirstParam) Data() (data interface{}, err error) {
	return p, nil
}
func (p FirstParam) Order() []exp.OrderedExpression {
	return sqlbuilder.ConcatOrderedExpression(goqu.C("Fid").Asc())
}
func (p FirstParam) Select() []any {
	return nil
}

func (p FirstParam) Where() (expressions sqlbuilder.Expressions, err error) {
	expressions = sqlbuilder.Expressions{
		goqu.C("Fid").Eq(p.Id),
		goqu.C("Fdeleted").Eq(""),
		goqu.C("Fname").ILike("%" + p.Name + "%"),
	}
	if p.Email != "" {
		expressions = append(expressions, goqu.C("Femail").Eq(p.Email))
	}
	return expressions, nil
}

func TestFirst(t *testing.T) {
	row := FirstParam{
		Id:   "123",
		Name: "张三",
	}
	sql, err := sqlbuilder.NewFirstBuilder(row).Merge(tenant.First(tenantData)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

package tenant_test

import (
	"fmt"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/tenant"
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

func (p InsertParam) Validate() (err error) {
	return nil
}

func (p InsertParam) Data() (data interface{}, err error) {
	return p, nil
}

func (p InsertParam) GetTenantField() tenant.TenantField {
	return tenant.TenantField{
		Name: "Ftenant",
		Value: func(in any) (value any, err error) {
			return p.Tenant, nil
		},
	}
}

var tenantData = tenant.TenantField{
	Name: "Ftenant",
	Value: func(in any) (value any, err error) {
		return "123", nil
	},
}

func TestInsert(t *testing.T) {
	row := InsertParam{
		Id:     "123",
		Name:   "张三",
		Email:  "莉丝",
		Tenant: "1000001",
	}
	sql, err := sqlbuilder.NewInsertBuilder(row).Merge(tenant.Insert(row)).ToSQL()
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

func (p UpdateParam) Where() (expressions []goqu.Expression, err error) {
	expressions = []goqu.Expression{
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

func (p FirstParam) Where() (expressions []goqu.Expression, err error) {
	expressions = []goqu.Expression{
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

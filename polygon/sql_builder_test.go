package polygon_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/polygon"
	"github.com/suifengpiao14/businessunit/tenant"
	"github.com/suifengpiao14/sqlbuilder"
)

func init() {
	sqlbuilder.Dialect = sqlbuilder.Dialect_Mysql
}

type InsertParam struct {
	Id     string `db:"Fid"`
	Name   string `db:"Fname"`
	Path   string `db:"Fpath"`
	Tenant string `db:"-"`
}

func (p InsertParam) Table() string {
	return "t_polygon"
}

type Polygon struct {
	Path string `json:"path"`
}

func (p Polygon) Points() (points polygon.Points, err error) {
	path := p.Path
	pointSet := make([][2]float64, 0)
	err = json.Unmarshal([]byte(path), &pointSet)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	geoPointSet := make(polygon.Points, 0)
	for _, point := range pointSet {
		geoPoint := polygon.Point{
			Lat: point[1],
			Lng: point[0],
		}
		geoPointSet = append(geoPointSet, geoPoint)
	}
	return geoPointSet, nil
}

func TestInsert(t *testing.T) {
	row := InsertParam{
		Id:     "123",
		Name:   "张三",
		Path:   `[[0,0],[2,0],[0,2],[2,2]]`,
		Tenant: "1000001",
	}
	polyg := Polygon{Path: row.Path}
	tenantField := tenant.NewTenant(row.Tenant).Field
	points, err := polyg.Points()
	require.NoError(t, err)
	boxField := points.GetBoundingBox()
	sql, err := sqlbuilder.NewInsertBuilder(row.Table()).AppendFields(tenantField).AppendFields(boxField.Fields()...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

type ListParam struct {
	Id        string `db:"Fid"`
	Name      string `db:"Fname"`
	Path      string `db:"Fpath"`
	Tenant    string `db:"-"`
	PageIndex string `db:"-"`
	PageSize  string `db:"-"`
}

func (p ListParam) Table() string {
	return "t_polygon"
}

func (p ListParam) Order() []exp.OrderedExpression {
	return sqlbuilder.ConcatOrderedExpression(goqu.C("Fid").Asc())
}

func (p ListParam) Pagination() (pageIndex int, pageSize int) {
	return cast.ToInt(p.PageIndex), cast.ToInt(p.PageSize)
}
func (p ListParam) Select() []any {
	return nil
}
func (p ListParam) Where() (expressions []exp.Expression, err error) {
	expressions = make([]exp.Expression, 0)
	if p.Name != "" {
		expressions = append(expressions, goqu.C("Fname").ILike("%"+p.Name+"%"))
	}
	return
}

func TestList(t *testing.T) {
	row := ListParam{
		Id:     "123",
		Name:   "张三",
		Path:   `[[0,0],[2,0],[0,2],[2,2]]`,
		Tenant: "1000001",
	}
	polyg := Polygon{Path: row.Path}
	points, err := polyg.Points()
	require.NoError(t, err)
	tenantField := tenant.NewTenant(row.Tenant).Field
	boxField := points.GetBoundingBox()
	sql, err := sqlbuilder.NewListBuilder(row.Table()).AppendFields(tenantField).AppendFields(boxField.Fields()...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

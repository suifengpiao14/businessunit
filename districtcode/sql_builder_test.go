package districtcode_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/districtcode"
	"github.com/suifengpiao14/sqlbuilder"
)

type District struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (d District) Table() sqlbuilder.TableConfig {
	return sqlbuilder.NewTableConfig("t_city_info")
}
func (d District) Select() (columns []any) {
	return
}
func (d District) Pagination() (index int, size int) {
	return
}

func TestGetChildren(t *testing.T) {
	d := District{
		Code: "440300",
		Name: "深圳",
	}
	codeField := sqlbuilder.NewField(func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return d.Code, nil }).SetName("code")
	nameField := sqlbuilder.NewField(func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return d.Name, nil }).SetName("name")
	districtcode.OptionsGetChildren(codeField, nameField, districtcode.Depth_max)
	sql, err := sqlbuilder.NewListBuilder(d.Table()).AppendFields(codeField, nameField).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

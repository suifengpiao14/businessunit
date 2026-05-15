package districtcode_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/districtcode"
	"gitlab.huishoubao.com/gopackage/sqlbuilder"
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
	codeField := sqlbuilder.NewField(d.Code).SetName("code")
	nameField := sqlbuilder.NewField(d.Name).SetName("name")
	districtcode.OptionsGetChildren(codeField, nameField, districtcode.Depth_max)
	fs := sqlbuilder.Fields{
		codeField,
		nameField,
	}
	sql, err := sqlbuilder.NewListBuilder(d.Table()).ToSQL(fs)
	require.NoError(t, err)
	fmt.Println(sql)
}

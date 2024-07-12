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

func (d District) GetDistrict() districtcode.District {
	filed := districtcode.NewDistrictField(
		func(in any) (any, error) { return d.Code, nil },
		func(in any) (any, error) { return d.Name, nil },
	)
	return districtcode.District{
		DistrictField: &filed,
	}
}
func (d District) Table() string {
	return "t_city_info"
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
	sql, err := sqlbuilder.NewListBuilder(d).Merge(districtcode.GetChildren(d, districtcode.Depth_max)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)
}

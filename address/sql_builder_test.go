package address_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/address"
	"github.com/suifengpiao14/sqlbuilder"
)

type InsertAddress struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	OwnerID      string `json:"ownerId"`
	Label        string `json:"label"`
	ContactPhone string `json:"contactPhone"`
	ContactName  string `json:"contactName"`
	Address      string `json:"address"`
	ProvinceId   string `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
	CityId       string `json:"cityId"`
	CityName     string `json:"cityName"`
	AreaId       string `json:"areaId"`
	AreaName     string `json:"areaName"`
	IsDefault    string `json:"isDefault"`
}

func (addr InsertAddress) GetAddress() (addres address.Address) {

	addres = address.Address{
		TenantID:     addr.TenantID,
		OwnerID:      addr.OwnerID,
		Label:        addr.Label,
		ContactPhone: addr.ContactPhone,
		ContactName:  addr.ContactName,
		Address:      addr.Address,
		IsDefault:    addr.IsDefault,
		ProvinceCode: addr.ProvinceId,
		ProvinceName: addr.ProvinceName,
		CityCode:     addr.CityId,
		CityName:     addr.CityName,
		AreaCode:     addr.AreaId,
		AreaName:     addr.AreaName,
	}
	addres.Init(addr.Table(), addr, addr)
	return addres
}

func (addr InsertAddress) Table() (table string) {
	return "t_address"
}

func (addr InsertAddress) CleanDefault(rawSql string) (err error) {
	fmt.Println(rawSql)
	return nil
}

func (addr InsertAddress) GetCount(rawSql string) (count int, err error) {
	fmt.Println(rawSql)
	return 0, nil
}

func TestInsertDoc(t *testing.T) {
	var addr = InsertAddress{}
	reqArgs := addr.GetAddress().Fields().DocRequestArgs()
	markdown := reqArgs.Makedown()
	fmt.Println(markdown)
	example := reqArgs.JsonExample(true)
	fmt.Println(example)
}

func TestInsertDDL(t *testing.T) {
	var addr = InsertAddress{}
	columns, err := addr.GetAddress().Fields().DBColumns()
	require.NoError(t, err)
	ddl := columns.DDL(sqlbuilder.Driver_mysql)
	fmt.Println(ddl)
}

func TestInsert(t *testing.T) {
	var addr = InsertAddress{
		OwnerID:      "123",
		Label:        "return",
		ContactPhone: "15999646794",
		ContactName:  "pz",
		Address:      "地球中国广东",
		ProvinceId:   "440000",
		ProvinceName: "广东",
		CityId:       "440300",
		CityName:     "深圳",
		AreaId:       "440301",
		AreaName:     "福田",
		IsDefault:    "1",
		TenantID:     "15",
	}
	sql, err := sqlbuilder.NewInsertBuilder(addr).AppendField(addr.GetAddress().Fields()...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

type UpdateAddress struct {
	InsertAddress
}

func TestUpdate(t *testing.T) {
	var addr = UpdateAddress{
		InsertAddress: InsertAddress{
			TenantID:     "478",
			OwnerID:      "124",
			Label:        "recive",
			ContactPhone: "15999646794",
			ContactName:  "pz",
			Address:      "地球中国广东",
			ProvinceId:   "440000",
			ProvinceName: "广东",
			CityId:       "440300",
			CityName:     "深圳",
			AreaId:       "440301",
			AreaName:     "福田",
			IsDefault:    "1",
		},
	}
	sql, err := sqlbuilder.NewUpdateBuilder(addr).AppendField(addr.GetAddress().Fields()...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

type ListAddress struct {
	InsertAddress
	PageIndex string `json:"pageIndex"`
	PageSize  string `json:"pageSize"`
}

func (l ListAddress) Pagination() (pageIndex int, pageSize int) {
	return cast.ToInt(l.PageIndex), cast.ToInt(pageSize)
}

func (l ListAddress) Select() []any {
	return nil
}

func TestSelect(t *testing.T) {
	var addr = ListAddress{
		InsertAddress: InsertAddress{
			TenantID:     "478",
			OwnerID:      "124",
			Label:        "recive",
			ContactPhone: "15999646794",
			ContactName:  "pz",
			Address:      "地球中国广东",
			ProvinceId:   "440000",
			ProvinceName: "广东",
			CityId:       "440300",
			CityName:     "深圳",
			AreaId:       "440301",
			AreaName:     "福田",
			IsDefault:    "1",
		},
	}
	sql, err := sqlbuilder.NewListBuilder(addr).AppendFields(addr.GetAddress().Fields()...).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

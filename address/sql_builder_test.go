package address_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/address"
	"github.com/suifengpiao14/businessunit/tenant"
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
		TenatID:      tenant.NewTenantField(func(in any) (any, error) { return addr.TenantID, nil }).AppendWhereFn(sqlbuilder.ValueFnDirect).SetName("businessId"),
		Label:        address.NewLabelField(func(in any) (any, error) { return addr.Label, nil }, nil),
		IsDefault:    address.NewIsDefaultField(func(in any) (any, error) { return addr.IsDefault, nil }, nil).AppendWhereFn(sqlbuilder.ValueFnDirect),
		ContactPhone: address.NewContactPhoneField(func(in any) (any, error) { return addr.ContactPhone, nil }).AppendWhereFn(sqlbuilder.ValueFnDirect),
		ContactName:  *address.NewContactNameField(func(in any) (any, error) { return addr.ContactName, nil }).SetName("contactName"),
		Address:      *address.NewAddressField(func(in any) (any, error) { return addr.Address, nil }),

		Province: address.NewProvinceField(func(in any) (any, error) { return addr.ProvinceName, nil }, func(in any) (any, error) { return addr.ProvinceId, nil }),
		City:     address.NewCityField(func(in any) (any, error) { return addr.CityId, nil }, func(in any) (any, error) { return addr.CityId, nil }),
		Area:     address.NewAreaField(func(in any) (any, error) { return addr.AreaId, nil }, func(in any) (any, error) { return addr.AreaName, nil }),
	}
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
func (addr InsertAddress) Fields() (fields sqlbuilder.Fields) {
	address := addr.GetAddress()
	fields = sqlbuilder.Fields{
		address.TenatID.Field,
		address.Label.Field,
		address.IsDefault.GetBooleanField().Field,
		address.ContactPhone.Field,
		address.ContactName,
		address.Address,
	}
	fields = append(fields, address.Province.GetIdTitle().Fields()...)
	fields = append(fields, address.City.GetIdTitle().Fields()...)
	fields = append(fields, address.Area.GetIdTitle().Fields()...)

	return fields
}

func TestInsertDoc(t *testing.T) {
	var addr = InsertAddress{}
	reqArgs, err := addr.Fields().DocRequestArgs()
	require.NoError(t, err)
	markdown := reqArgs.Makedown()
	fmt.Println(markdown)
	example := reqArgs.JsonExample(true)
	fmt.Println(example)
}

func TestInsertDDL(t *testing.T) {
	var addr = InsertAddress{}
	columns, err := addr.Fields().DBColumns()
	require.NoError(t, err)
	ddl := columns.DDL(sqlbuilder.Dialect_mysql)
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
	sql, err := sqlbuilder.NewInsertBuilder(addr).Merge(address.Insert(addr, addr, addr)).ToSQL()
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
	sql, err := sqlbuilder.NewUpdateBuilder(addr).Merge(address.Update(addr, addr)).ToSQL()
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
	sql, err := sqlbuilder.NewListBuilder(addr).Merge(address.List(addr)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

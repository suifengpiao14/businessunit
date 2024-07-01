package address_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/address"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

type Address struct {
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
}

func (addr Address) GetAddress() (addres address.Address) {
	addres = address.Address{
		OwnerID: sqlbuilder.Field{
			Name: "Fowner_id",
			Value: func(in any) (value any, err error) {
				return addr.OwnerID, nil
			},
			WhereValue: func(in any) (value any, err error) {
				return addr.OwnerID, nil
			},
		},
		Label: sqlbuilder.Field{
			Name: "Flabel",
			Value: func(in any) (value any, err error) {
				return addr.Label, nil
			},
		},
		ContactPhone: sqlbuilder.Field{
			Name: "Fcontact_phone",
			Value: func(in any) (value any, err error) {
				return addr.ContactPhone, nil
			},
		},
		ContactName: sqlbuilder.Field{
			Name: "Fcontact_name",
			Value: func(in any) (value any, err error) {
				return addr.ContactName, nil
			},
		},
		Address: sqlbuilder.Field{
			Name: "Faddress",
			Value: func(in any) (value any, err error) {
				return addr.Address, nil
			},
		},

		Province: title.Title{
			ID: identity.IdentityField{
				Name: "Fprovice_id",
				Value: func(in any) (value any, err error) {
					return addr.ProvinceId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fprovice",
				Value: func(in any) (value any, err error) {
					return addr.ProvinceName, nil
				},
			},
		},
		City: title.Title{
			ID: identity.IdentityField{
				Name: "Fcity_id",
				Value: func(in any) (value any, err error) {
					return addr.CityId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fcity",
				Value: func(in any) (value any, err error) {
					return addr.CityName, nil
				},
			},
		},
		Area: title.Title{
			ID: identity.IdentityField{
				Name: "Farea_id",
				Value: func(in any) (value any, err error) {
					return addr.AreaId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Farea",
				Value: func(in any) (value any, err error) {
					return addr.AreaName, nil
				},
			},
		},
	}
	return addres
}

func (addr Address) Table() (table string) {
	return "t_address"
}

var addr = Address{
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
}

func TestInsert(t *testing.T) {
	sql, err := sqlbuilder.NewInsertBuilder(addr).Merge(address.Insert(addr)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

func TestUpdate(t *testing.T) {
	sql, err := sqlbuilder.NewUpdateBuilder(addr).Merge(address.Update(addr)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

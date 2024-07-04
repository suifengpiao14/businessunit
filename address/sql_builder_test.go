package address_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/address"
	"github.com/suifengpiao14/businessunit/boolean"
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/businessunit/phone"
	"github.com/suifengpiao14/businessunit/tenant"
	"github.com/suifengpiao14/businessunit/title"
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
		TenatID: tenant.TenantField{
			Name: "Fbusiness_id",
			ValueFn: func(in any) (value any, err error) {
				return addr.TenantID, nil
			},
			WhereValueFn: func(in any) (value any, err error) {
				return addr.TenantID, nil
			},
		},
		OwnerID: ownerid.OwnerIdField{
			Name: "Fowner_id",
			ValueFn: func(in any) (value any, err error) {
				return addr.OwnerID, nil
			},
			WhereValueFn: func(in any) (value any, err error) {
				return addr.OwnerID, nil
			},
		},
		Label: enum.EnumField{
			Field: sqlbuilder.Field{
				Name: "Flabel",
				ValueFn: func(in any) (value any, err error) {
					return addr.Label, nil
				},
				WhereValueFn: func(in any) (value any, err error) {
					return addr.Label, nil
				},
			},
			EnumTitles: enum.EnumTitles{
				{
					Key:   "recive",
					Title: "收获地址",
				},
				{
					Key:   "return",
					Title: "退货地址",
				},
			},
		},
		IsDefault: boolean.BooleanField{
			Field: sqlbuilder.Field{
				Name: "Fis_default",
				ValueFn: func(in any) (value any, err error) {
					return addr.IsDefault, nil
				},
			},
			TrueFalseTitleFn: func() (trueTitle enum.EnumTitle, falseTitle enum.EnumTitle) {
				trueTitle = enum.EnumTitle{
					Key:   "1",
					Title: "是",
				}
				falseTitle = enum.EnumTitle{
					Key:   "2",
					Title: "否",
				}
				return trueTitle, falseTitle
			},
		},
		ContactPhone: phone.PhoneField{
			Name: "Fcontact_phone",
			ValueFn: func(in any) (value any, err error) {
				return addr.ContactPhone, nil
			},
			WhereValueFn: func(in any) (value any, err error) {
				return addr.ContactPhone, nil
			},
		},
		ContactName: sqlbuilder.Field{
			Name: "Fcontact_name",
			ValueFn: func(in any) (value any, err error) {
				return addr.ContactName, nil
			},
			WhereValueFn: func(in any) (value any, err error) {
				return sqlbuilder.Ilike{"%", addr.ContactName}, nil
			},
		},
		Address: sqlbuilder.Field{
			Name: "Faddress",
			ValueFn: func(in any) (value any, err error) {
				return addr.Address, nil
			},
		},

		Province: title.Title{
			ID: identity.IdentityField{
				Name: "Fprovice_id",
				ValueFn: func(in any) (value any, err error) {
					return addr.ProvinceId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fprovice",
				ValueFn: func(in any) (value any, err error) {
					return addr.ProvinceName, nil
				},
			},
		},
		City: title.Title{
			ID: identity.IdentityField{
				Name: "Fcity_id",
				ValueFn: func(in any) (value any, err error) {
					return addr.CityId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fcity",
				ValueFn: func(in any) (value any, err error) {
					return addr.CityName, nil
				},
			},
		},
		Area: title.Title{
			ID: identity.IdentityField{
				Name: "Farea_id",
				ValueFn: func(in any) (value any, err error) {
					return addr.AreaId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Farea",
				ValueFn: func(in any) (value any, err error) {
					return addr.AreaName, nil
				},
			},
		},
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
	}
	sql, err := sqlbuilder.NewInsertBuilder(addr).Merge(address.Insert(addr)).ToSQL()
	require.NoError(t, err)
	fmt.Println(sql)

}

type UpdateAddress struct {
	InsertAddress
}

func (addr UpdateAddress) Table() (table string) {
	return "t_address"
}

func (addr UpdateAddress) CleanDefault(rawSql string) (err error) {
	fmt.Println(rawSql)
	return nil
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
	sql, err := sqlbuilder.NewUpdateBuilder(addr).Merge(address.Update(addr)).ToSQL()
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

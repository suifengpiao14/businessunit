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
			Field: sqlbuilder.Field{
				Name: "Fbusiness_id",
				ValueFns: func(in any) (value any, err error) {
					return addr.TenantID, nil
				},
				WhereFormatFns: sqlbuilder.FormatFns{sqlbuilder.DirectFormat},
			},
		},
		// OwnerID: ownerid.NewOwnerIdField("Fowner_id", func(in any) (value any, err error) {
		// 	return addr.OwnerID, nil
		// },
		// 	func(in any) (value any, err error) {
		// 		return addr.OwnerID, nil
		// 	},
		// 	nil,
		// ),
		Label: enum.EnumField{
			Field: sqlbuilder.Field{
				Name: "Flabel",
				ValueFns: func(in any) (value any, err error) {
					return addr.Label, nil
				},
				WhereFormatFns: sqlbuilder.FormatFns{sqlbuilder.DirectFormat},
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
				ValueFns: func(in any) (value any, err error) {
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
			ValueFns: func(in any) (value any, err error) {
				return addr.ContactPhone, nil
			},
			WhereFormatFns: sqlbuilder.FormatFns{sqlbuilder.DirectFormat},
		},
		ContactName: sqlbuilder.Field{
			Name: "Fcontact_name",
			ValueFns: func(in any) (value any, err error) {
				return addr.ContactName, nil
			},
			WhereFormatFns: sqlbuilder.FormatFns{sqlbuilder.DirectFormat},
		},
		Address: sqlbuilder.Field{
			Name: "Faddress",
			ValueFns: func(in any) (value any, err error) {
				return addr.Address, nil
			},
		},

		Province: title.Title{
			ID: identity.IdentityField{
				Name: "Fprovice_id",
				ValueFns: func(in any) (value any, err error) {
					return addr.ProvinceId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fprovice",
				ValueFns: func(in any) (value any, err error) {
					return addr.ProvinceName, nil
				},
			},
		},
		City: title.Title{
			ID: identity.IdentityField{
				Name: "Fcity_id",
				ValueFns: func(in any) (value any, err error) {
					return addr.CityId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Fcity",
				ValueFns: func(in any) (value any, err error) {
					return addr.CityName, nil
				},
			},
		},
		Area: title.Title{
			ID: identity.IdentityField{
				Name: "Farea_id",
				ValueFns: func(in any) (value any, err error) {
					return addr.AreaId, nil
				},
			},
			Title: sqlbuilder.Field{
				Name: "Farea",
				ValueFns: func(in any) (value any, err error) {
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

func (addr InsertAddress) GetCount(rawSql string) (count int, err error) {
	fmt.Println(rawSql)
	return 0, nil
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

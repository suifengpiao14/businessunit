package address

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

var Time_format = sqlbuilder.Time_format

/**
 -- 多业务，业务组合时增加 tenant 即可
CREATE TABLE `t_merchant_address_info` (
  `Fid` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `Fowner_id` int(11) NOT NULL COMMENT '地址所有者ID',
  `Flabel` enum('returnAddress','deliveryAddress') NOT NULL COMMENT '标签(returnAddress-退货地址,deliveryAddress-收货地址)',
  `Fcontact_phone` varchar(11) NOT NULL DEFAULT '' COMMENT '联系电话',
  `Fcontact_name` varchar(32) NOT NULL DEFAULT '' COMMENT '联系人',
  `Fprovince_id` int(11) NOT NULL DEFAULT '0' COMMENT '收货地址省id',
  `Fprovince_name` varchar(32) NOT NULL DEFAULT '' COMMENT '收货地址省名称',
  `Fcity_id` int(11) NOT NULL DEFAULT '0' COMMENT '收货地址市id',
  `Fcity_name` varchar(32) NOT NULL DEFAULT '' COMMENT '收货地址市名称',
  `Farea_id` int(11) NOT NULL DEFAULT '0' COMMENT '收货地址区id',
  `Farea_name` varchar(32) NOT NULL DEFAULT '' COMMENT '收货地址区名称',
  `Faddress` varchar(128) NOT NULL DEFAULT '' COMMENT '收货地址',
  `Faddress_all` varchar(256) NOT NULL DEFAULT '' COMMENT '收货地址(完整地址)',
  `Fis_default` enum('1','2') DEFAULT '2' COMMENT '默认地址1默认,2-非默认',
  `Fdelete_at` datetime DEFAULT NULL COMMENT '删除时间',
  `Fauto_create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `Fauto_update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`Fid`),
  KEY `Fowner_id` (`Fowner_id`,`Flabel`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COMMENT='地址信息表';

**/

type Address struct {
	OwnerID      sqlbuilder.Field
	Label        sqlbuilder.Field
	ContactPhone sqlbuilder.Field
	ContactName  sqlbuilder.Field
	Address      sqlbuilder.Field

	Province title.TitleI
	City     title.TitleI
	Area     title.TitleI
}

func (address Address) Fields() (fileds sqlbuilder.Fields) {
	fileds = make(sqlbuilder.Fields, 0)
	fileds = append(fileds,
		address.OwnerID,
		address.Label,
		address.ContactPhone,
		address.ContactName,
		address.Address,
	)
	return fileds
}

type AddressI interface {
	GetAddress() Address
}

func _DataFn(addressI AddressI) sqlbuilder.DataFn {
	return func() (any, error) {
		address := addressI.GetAddress()
		m, err := address.Fields().Map()
		return m, err
	}
}

func _WhereFn(addressI AddressI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		address := addressI.GetAddress()
		expressions = make([]goqu.Expression, 0)
		for _, field := range address.Fields() {
			if field.WhereValue == nil {
				continue
			}
			val, err := field.WhereValue(nil)
			if err != nil {
				return nil, err
			}
			if ex, ok := sqlbuilder.TryParseExpressions(field.Name, val); ok {
				expressions = append(expressions, ex...)
			}
			expressions = append(expressions, goqu.C(field.Name).Eq(val))
		}
		return expressions, nil
	}
}

func Insert(addressI AddressI) sqlbuilder.InsertParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(addressI)).Merge(
		title.Insert(provice),
		title.Insert(city),
		title.Insert(area),
	)
}

func Update(addressI AddressI) sqlbuilder.UpdateParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(addressI)).AppendWhere(_WhereFn(addressI)).Merge(
		title.Update(provice),
		title.Update(city),
		title.Update(area),
	)
}

func First(addressI AddressI) sqlbuilder.FirstParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.First(provice),
		title.First(city),
		title.First(area),
	)
}

func List(addressI AddressI) sqlbuilder.ListParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.List(provice),
		title.List(city),
		title.List(area),
	)
}

func Total(addressI AddressI) sqlbuilder.TotalParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.Total(provice),
		title.Total(city),
		title.Total(area),
	)
}

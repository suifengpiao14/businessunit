package address

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/boolean"
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/businessunit/phone"
	"github.com/suifengpiao14/businessunit/tenant"
	"github.com/suifengpiao14/businessunit/title"
	"github.com/suifengpiao14/sqlbuilder"
)

type AddressRule struct {
	TenatID   tenant.TenantField // 业务、应用、租户等唯一标识
	OwnerID   ownerid.OwnerIdField
	Label     enum.EnumField
	MaxNumber sqlbuilder.Field // 单个业务下指定类型可配置最大条数
}

type AddressRules []AddressRule

func (rs AddressRules) GetByLabel(tenatID tenant.TenantField, ownerID ownerid.OwnerIdField, label enum.EnumField) (addressRule *AddressRule, exist bool) {

}

// GetAddressRules 业务方重新初始化该函数，获取值
var GetAddressRules = func() (addressRules AddressRules) {
	return
}

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
	TenatID      tenant.TenantField // 业务、应用、租户等唯一标识
	OwnerID      ownerid.OwnerIdField
	Label        enum.EnumField
	ContactPhone phone.PhoneField
	ContactName  sqlbuilder.Field
	Address      sqlbuilder.Field
	IsDefault    boolean.BooleanI

	Province title.TitleI
	City     title.TitleI
	Area     title.TitleI
}

func (address Address) Fields() (fileds sqlbuilder.Fields) {
	fileds = make(sqlbuilder.Fields, 0)
	fileds = append(fileds,
		address.ContactName,
		address.Address,
	)
	return fileds
}

type AddressI interface {
	GetAddress() Address
	sqlbuilder.Table
	CleanDefault(rawSql string) (err error) // 清除默认记录标记(新写数据有默认属性)
}

func _DataFn(addressI AddressI) sqlbuilder.DataFn {
	return func() (any, error) {
		address := addressI.GetAddress()
		isDefaultField := address.IsDefault
		if isDefaultField != nil && isDefaultField.GetBooleanField().IsTrue() { // 设置为默认记录前，清除数据库当前默认标记
			// 构造一个false 值记录
			falseField := boolean.Switch(isDefaultField)
			rawSql, err := sqlbuilder.NewUpdateBuilder(addressI).Merge(
				boolean.Update(falseField),
			).AppendWhere(
				tenant.WhereFn(address.TenatID),
				ownerid.WhereFn(address.OwnerID),
				enum.WhereFn(address.Label),
			).ToSQL()
			if err != nil {
				return nil, err
			}
			err = addressI.CleanDefault(rawSql)
			if err != nil {
				return nil, err
			}
		}
		m, err := address.Fields().Map()
		return m, err
	}
}

func _WhereFn(addressI AddressI) sqlbuilder.WhereFn {
	return func() (expressions []goqu.Expression, err error) {
		address := addressI.GetAddress()
		validate := validator.New()
		err = validate.Struct(address)
		if err != nil {
			return nil, err
		}
		ownerID, err := address.OwnerID.WhereValueFn(nil)
		if err != nil {
			return nil, err
		}
		if cast.ToString(ownerID) == "" {
			return nil, errors.Errorf("字段%s不能为空", address.OwnerID.Name)
		}
		expressions = make([]goqu.Expression, 0)
		for _, field := range address.Fields() {
			if field.WhereValueFn == nil {
				continue
			}
			val, err := field.WhereValueFn(nil)
			if err != nil {
				return nil, err
			}
			if ex, ok := sqlbuilder.TryParseExpressions(field.Name, val); ok {
				expressions = append(expressions, ex...)
				continue
			}
			expressions = append(expressions, goqu.Ex{field.Name: val})
		}
		return expressions, nil
	}
}

func _OrderFn(booleanI boolean.BooleanI) sqlbuilder.OrderFn { // 默认记录排前面
	return func() (orderedExpressions []exp.OrderedExpression) {
		field := booleanI.GetBooleanField()
		trueTitle, falseTitle := booleanI.GetTrueFalseTitle()
		segment := fmt.Sprintf("FIELD(`%s`, ?, ?)", field.Name)
		expression := goqu.L(segment, trueTitle.Key, falseTitle.Key)
		orderedExpression := exp.NewOrderedExpression(expression, exp.AscDir, exp.NoNullsSortType)
		orderedExpressions = sqlbuilder.ConcatOrderedExpression(orderedExpression)
		return orderedExpressions
	}
}

func Insert(addressI AddressI) sqlbuilder.InsertParam {
	address := addressI.GetAddress()
	phoneField := address.ContactPhone
	provice := address.Province
	city := address.City
	area := address.Area
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(addressI)).Merge(
		title.Insert(provice),
		title.Insert(city),
		title.Insert(area),
		phone.Insert(phoneField),
		enum.Insert(label),
		boolean.Insert(isDefault),
		tenant.Insert(address.TenatID),
		ownerid.Insert(address.OwnerID),
	)
}

func Update(addressI AddressI) sqlbuilder.UpdateParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	phoneField := address.ContactPhone
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(addressI)).AppendWhere(_WhereFn(addressI)).Merge(
		title.Update(provice),
		title.Update(city),
		title.Update(area),
		phone.Update(phoneField),
		enum.Update(label),
		boolean.Update(isDefault),
		tenant.Update(address.TenatID),
		ownerid.Update(address.OwnerID),
	)
}

func First(addressI AddressI) sqlbuilder.FirstParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	phoneField := address.ContactPhone
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.First(provice),
		title.First(city),
		title.First(area),
		phone.First(phoneField),
		enum.First(label),
		boolean.First(isDefault),
		tenant.First(address.TenatID),
		ownerid.First(address.OwnerID),
	).AppendOrder(_OrderFn(isDefault))
}

func List(addressI AddressI) sqlbuilder.ListParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	phoneField := address.ContactPhone
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.List(provice),
		title.List(city),
		title.List(area),
		phone.List(phoneField),
		enum.List(label),
		boolean.List(isDefault),
		tenant.List(address.TenatID),
		ownerid.List(address.OwnerID),
	).AppendOrder(_OrderFn(isDefault))
}

func Total(addressI AddressI) sqlbuilder.TotalParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	phoneField := address.ContactPhone
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		title.Total(provice),
		title.Total(city),
		title.Total(area),
		phone.Total(phoneField),
		enum.Total(label),
		boolean.Total(isDefault),
		tenant.Total(address.TenatID),
		ownerid.Total(address.OwnerID),
	)
}

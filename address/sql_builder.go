package address

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/boolean"
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/businessunit/idtitle"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/businessunit/phone"
	"github.com/suifengpiao14/businessunit/tenant"
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
	for _, r := range rs {
		if r.TenatID.IsEqual(tenatID) && r.OwnerID.IsEqual(ownerID.Field) && r.Label.IsEqual(label.Field) {
			return &r, true
		}
	}
	return nil, false
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

	Province idtitle.IdTitleI
	City     idtitle.IdTitleI
	Area     idtitle.IdTitleI
}

func (address Address) Fields() (fileds sqlbuilder.Fields) {
	fileds = make(sqlbuilder.Fields, 0)
	fileds = append(fileds,
		address.ContactName,
		address.Address,
	)
	return fileds
}

func NewContactNameField(valueFn sqlbuilder.ValueFn) sqlbuilder.Field {
	return sqlbuilder.NewField(valueFn).SetName("contact_name").SetTitle("联系人")
}
func NewContactPhoneField(valueFn sqlbuilder.ValueFn) phone.PhoneField {
	return phone.NewPhoneField(valueFn).SetName("contact_phone").SetTitle("联系手机号")
}

func NewAddressField(valueFn sqlbuilder.ValueFn) sqlbuilder.Field {
	return sqlbuilder.NewField(valueFn).SetName("address").SetTitle("详细地址")
}

// NewProvinceField 封装省字段
func NewProvinceField(idValueFn sqlbuilder.ValueFn, titleValueFn sqlbuilder.ValueFn) idtitle.IdTitleI {
	idTitle := idtitle.NewIdTitleFiled(idValueFn, titleValueFn)
	idTitle.ID.SetName("provice_id").SetTitle("省ID")
	idTitle.Title.SetName("provice").SetTitle("省")
	return idTitle
}

// NewCityField 封装市字段
func NewCityField(idValueFn sqlbuilder.ValueFn, titleValueFn sqlbuilder.ValueFn) idtitle.IdTitleI {
	idTitle := idtitle.NewIdTitleFiled(idValueFn, titleValueFn)
	idTitle.ID.SetName("city_id").SetTitle("城市ID")
	idTitle.Title.SetName("city").SetTitle("城市")
	return idTitle
}

// NewAreaField 封装区字段
func NewAreaField(idValueFn sqlbuilder.ValueFn, titleValueFn sqlbuilder.ValueFn) idtitle.IdTitleI {
	idTitle := idtitle.NewIdTitleFiled(idValueFn, titleValueFn)
	idTitle.ID.SetName("area_id").SetTitle("区ID")
	idTitle.Title.SetName("area").SetTitle("区")
	return idTitle
}

type AddressI interface {
	GetAddress() Address
	sqlbuilder.Table
}

// WithDefaultI 需要设置默认地址时,需要实现该接口
type WithDefaultI interface {
	CleanDefault(rawSql string) (err error)
}

type CheckRuleI interface {
	GetCount(rawSql string) (count int, err error) // 某种类型需要限制数量时,需要实现该接口,查询数据库已有的数量
}

func _DataFn(addressI AddressI) sqlbuilder.DataFn {
	return addressI.GetAddress().Fields().Data
}

func _WhereFn(addressI AddressI) sqlbuilder.WhereFn {
	return addressI.GetAddress().Fields().Where
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

func _ValidateRuleFn(addressI AddressI, checkRuleI CheckRuleI) sqlbuilder.ValidateFn {
	return func(_ any) error {
		if checkRuleI == nil {
			return nil
		}
		address := addressI.GetAddress()
		r, ok := GetAddressRules().GetByLabel(address.TenatID, address.OwnerID, address.Label)
		if !ok {
			return nil
		}
		val, err := r.MaxNumber.GetValue(nil)
		if err != nil {
			return err
		}
		maxNumber := cast.ToInt(val)
		if maxNumber == 0 {
			return nil
		}
		rawSql, err := sqlbuilder.NewTotalBuilder(addressI).AppendWhere(
			tenant.WhereFn(address.TenatID),
			ownerid.WhereFn(address.OwnerID),
			enum.WhereFn(address.Label),
		).ToSQL()
		if err != nil {
			return err
		}
		count, err := checkRuleI.GetCount(rawSql)
		if err != nil {
			return err
		}
		if maxNumber >= count {
			err = errors.Errorf(
				"%s-%s-%s-已有数量(%d)-超过最大数量限制(%d)",
				address.TenatID.LogString(),
				address.OwnerID.LogString(),
				address.Label.LogString(),
				count,
				maxNumber,
			)
			return err
		}

		return nil
	}
}

// _DealDefault 当前记录需要设置为默认记录时,清除已有的默认记录
func _DealDefault(addressI AddressI, withDWithDefaultI WithDefaultI) sqlbuilder.ValidateFn {
	return func(val any) (err error) {
		if withDWithDefaultI == nil {
			return nil
		}

		address := addressI.GetAddress()
		isDefaultField := address.IsDefault
		if isDefaultField == nil || !isDefaultField.GetBooleanField().IsTrue() { //当前记录不是默认记录时,无需处理
			return nil
		}

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
			return err
		}
		err = withDWithDefaultI.CleanDefault(rawSql)
		if err != nil {
			return err
		}
		return nil
	}
}

func Insert(addressI AddressI, withDefaultI WithDefaultI, validateRuleI CheckRuleI) sqlbuilder.InsertParam {
	address := addressI.GetAddress()
	phoneField := address.ContactPhone
	provice := address.Province
	city := address.City
	area := address.Area
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(addressI)).Merge(
		idtitle.Insert(provice),
		idtitle.Insert(city),
		idtitle.Insert(area),
		phone.Insert(phoneField),
		enum.Insert(label),
		boolean.Insert(isDefault),
		tenant.Insert(address.TenatID),
		ownerid.Insert(address.OwnerID),
	).AppendValidate(_ValidateRuleFn(addressI, validateRuleI), _DealDefault(addressI, withDefaultI))
}

func Update(addressI AddressI, withDefaultI WithDefaultI) sqlbuilder.UpdateParam {
	address := addressI.GetAddress()
	provice := address.Province
	city := address.City
	area := address.Area
	phoneField := address.ContactPhone
	label := address.Label
	isDefault := address.IsDefault
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(addressI)).AppendWhere(_WhereFn(addressI)).Merge(
		idtitle.Update(provice),
		idtitle.Update(city),
		idtitle.Update(area),
		phone.Update(phoneField),
		enum.Update(label),
		boolean.Update(isDefault),
		tenant.Update(address.TenatID),
		ownerid.Update(address.OwnerID),
	).AppendValidate(_DealDefault(addressI, withDefaultI))
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
		idtitle.First(provice),
		idtitle.First(city),
		idtitle.First(area),
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
		idtitle.List(provice),
		idtitle.List(city),
		idtitle.List(area),
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
		idtitle.Total(provice),
		idtitle.Total(city),
		idtitle.Total(area),
		phone.Total(phoneField),
		enum.Total(label),
		boolean.Total(isDefault),
		tenant.Total(address.TenatID),
		ownerid.Total(address.OwnerID),
	)
}

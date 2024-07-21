package address

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/boolean"
	"github.com/suifengpiao14/businessunit/districtcode"
	"github.com/suifengpiao14/businessunit/enum"
	"github.com/suifengpiao14/businessunit/keytitle"
	"github.com/suifengpiao14/businessunit/name"
	"github.com/suifengpiao14/businessunit/ownerid"
	"github.com/suifengpiao14/businessunit/phone"
	"github.com/suifengpiao14/businessunit/tenant"
	"github.com/suifengpiao14/sqlbuilder"
)

type AddressRule struct {
	TenatID   *tenant.Tenant // 业务、应用、租户等唯一标识
	OwnerID   *ownerid.OwnerID
	Label     *enum.Enum
	MaxNumber sqlbuilder.Field // 单个业务下指定类型可配置最大条数
}

type AddressRules []AddressRule

func (rs AddressRules) GetByLabel(tenatID *tenant.Tenant, ownerID *ownerid.OwnerID, label *enum.Enum) (addressRule *AddressRule, exist bool) {
	for _, r := range rs {
		if r.TenatID.Field.IsEqual(*tenatID.Field) && r.OwnerID.Field.IsEqual(*ownerID.Field) && r.Label.Field.IsEqual(*label.Field) {
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
	table        string
	TenantID     string
	OwnerID      string
	Label        string
	ContactPhone string
	ContactName  string
	Address      string
	IsDefault    string
	ProvinceCode string
	ProvinceName string
	CityCode     string
	CityName     string
	AreaCode     string
	AreaName     string

	TenatIDField      *tenant.Tenant // 业务、应用、租户等唯一标识
	OwnerIDField      *ownerid.OwnerID
	LabelField        *enum.Enum
	ContactPhoneField *phone.Phone
	ContactNameField  *name.Name
	AddressField      *sqlbuilder.Field
	IsDefaultField    *boolean.Boolean

	ProvinceField *districtcode.District
	CityField     *districtcode.District
	AreaField     *districtcode.District
}

func (addr *Address) Init(table string, checkRuleI CheckRuleI, withDWithDefaultI WithDefaultI) {
	addr.TenatIDField = tenant.NewTenant(addr.TenantID)
	addr.OwnerIDField = ownerid.NewOwnerID(addr.OwnerID)
	addr.LabelField = enum.NewEnum(addr.Label, sqlbuilder.Enums{
		{
			Key:       "recive",
			Title:     "收获地址",
			IsDefault: true,
		},
		{
			Key:   "return",
			Title: "退货地址",
		},
	})
	addr.LabelField.Field.SetName("label").SetTitle("标签")

	addr.ContactPhoneField = phone.NewPhone(addr.ContactPhone).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("contactPhone").SetTitle("联系手机号")
	})
	addr.ContactNameField = name.NewName(addr.ContactName)
	addr.ContactNameField.Field.SetName("contactName").SetTitle("联系人").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		MinLength: 1,
		MaxLength: 20, // 常规名称在20个字以内
		Type:      sqlbuilder.Schema_Type_string,
	})
	addr.AddressField = sqlbuilder.NewField(func(in any) (any, error) { return addr.Address, nil }).SetName("adrees").SetTitle("详细地址").MergeSchema(sqlbuilder.Schema{
		MaxLength: 100, // 线上统计最大55个字符，设置100 应该适合大部分场景大小
	})
	addr.IsDefaultField = boolean.NewBoolean(addr.IsDefault,
		sqlbuilder.Enum{
			Key:   "1",
			Title: "是",
		},
		sqlbuilder.Enum{
			Key:   "2",
			Title: "否",
		},
	)
	addr.IsDefaultField.Field.SetName("isDefault").SetTitle("默认")

	addr.ProvinceField = districtcode.NewDistrict(addr.ProvinceCode, addr.ProvinceName)
	addr.ProvinceField.CodeField.SetName("proviceId").SetTitle("省ID")
	addr.ProvinceField.NameField.SetName("provice").SetTitle("省").MergeSchema(sqlbuilder.Schema{
		MaxLength: 16,
	})

	addr.CityField = districtcode.NewDistrict(addr.ProvinceCode, addr.ProvinceName)
	addr.CityField.CodeField.SetName("cityId").SetTitle("城市ID")
	addr.CityField.NameField.SetName("city").SetTitle("城市").MergeSchema(sqlbuilder.Schema{
		MaxLength: 32, //海南省黎母山林场（海南黎母山省级自然保护区管理站）线上最长 25，设置32 2的倍数
	})

	addr.AreaField = districtcode.NewDistrict(addr.AreaCode, addr.AreaName)
	addr.AreaField.CodeField.SetName("areaId").SetTitle("区ID")
	addr.AreaField.NameField.SetName("area").SetTitle("区").MergeSchema(sqlbuilder.Schema{
		MaxLength: 128, //海南省黎母山林场（海南黎母山省级自然保护区管理站）线上最长 25，设置32 2的倍数
	})

	fields := addr.Fields()
	first := fields[0]
	first.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(
			_ValidateRuleFn(table, *addr, checkRuleI),
			_DealDefault(table, *addr, withDWithDefaultI),
		)
	})

}

func (addr Address) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		addr.TenatIDField.Field,
		addr.OwnerIDField.Field,
		addr.LabelField.Field,
		addr.ContactNameField.Field,
		addr.ContactPhoneField.Field,
		addr.AddressField,
		addr.IsDefaultField.Field,
	}
	fs.Append(addr.ProvinceField.Fields()...)
	fs.Append(addr.CityField.Fields()...)
	fs.Append(addr.AreaField.Fields()...)
	return fs
}

type AddressI interface {
	GetAddress() Address
	sqlbuilder.TableI
}

// WithDefaultI 需要设置默认地址时,需要实现该接口
type WithDefaultI interface {
	CleanDefault(rawSql string) (err error)
}

type CheckRuleI interface {
	GetCount(rawSql string) (count int, err error) // 某种类型需要限制数量时,需要实现该接口,查询数据库已有的数量
}

func _OrderFn(boolean boolean.Boolean) sqlbuilder.OrderFn { // 默认记录排前面
	return func() (orderedExpressions []exp.OrderedExpression) {
		segment := fmt.Sprintf("FIELD(`%s`, ?, ?)", boolean.Field.DBName())
		expression := goqu.L(segment, boolean.TrueEnum.Key, boolean.FalseEnum.Key)
		orderedExpression := exp.NewOrderedExpression(expression, exp.AscDir, exp.NoNullsSortType)
		orderedExpressions = sqlbuilder.ConcatOrderedExpression(orderedExpression)
		return orderedExpressions
	}
}

func _ValidateRuleFn(table string, address Address, checkRuleI CheckRuleI) sqlbuilder.ValueFn {
	return func(in any) (any, error) {
		if checkRuleI == nil {
			return in, nil
		}
		r, ok := GetAddressRules().GetByLabel(address.TenatIDField, address.OwnerIDField, address.LabelField)
		if !ok {
			return in, nil
		}
		val, err := r.MaxNumber.GetValue()
		if err != nil {
			return nil, err
		}
		maxNumber := cast.ToInt(val)
		if maxNumber == 0 {
			return in, nil
		}
		var tableFn sqlbuilder.TableFn = func() (table string) { return table }
		rawSql, err := sqlbuilder.NewTotalBuilder(tableFn).AppendFields(
			address.TenatIDField.Field,
			address.OwnerIDField.Field,
			address.LabelField.Field,
		).ToSQL()
		if err != nil {
			return nil, err
		}
		count, err := checkRuleI.GetCount(rawSql)
		if err != nil {
			return nil, err
		}
		if maxNumber >= count {
			err = errors.Errorf(
				"%s-%s-%s-已有数量(%d)-超过最大数量限制(%d)",
				address.TenatIDField.Field.LogString(),
				address.OwnerIDField.Field.LogString(),
				address.LabelField.Field.LogString(),
				count,
				maxNumber,
			)
			return nil, err
		}

		return in, nil
	}
}

// _DealDefault 当前记录需要设置为默认记录时,清除已有的默认记录
func _DealDefault(table string, address Address, withDWithDefaultI WithDefaultI) sqlbuilder.ValueFn {
	return func(val any) (any, error) {
		if withDWithDefaultI == nil {
			return val, nil
		}

		isDefaultField := address.IsDefaultField
		if isDefaultField == nil || !isDefaultField.IsTrue() { //当前记录不是默认记录时,无需处理
			return val, nil
		}

		// 构造一个false 值记录
		falseField := boolean.Switch(*isDefaultField)
		labelField := address.LabelField.Field.Copy()
		labelField.ShieldUpdate(true)
		var tableFn sqlbuilder.TableFn = func() (table string) { return table }
		rawSql, err := sqlbuilder.NewUpdateBuilder(tableFn).AppendField(falseField.Field).AppendField(
			address.TenatIDField.Field,
			address.OwnerIDField.Field,
			labelField,
		).ToSQL()
		if err != nil {
			return nil, err
		}
		err = withDWithDefaultI.CleanDefault(rawSql)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
}

func Insert(addressI AddressI, withDefaultI WithDefaultI, validateRuleI CheckRuleI) sqlbuilder.InsertParam {
	address := addressI.GetAddress()
	return sqlbuilder.NewInsertBuilder(nil).AppendField(address.Fields()...)
}

func Update(addressI AddressI, withDefaultI WithDefaultI) sqlbuilder.UpdateParam {
	address := addressI.GetAddress()
	provice := address.ProvinceField
	city := address.CityField
	area := address.AreaField
	phoneField := address.ContactPhoneField
	label := address.LabelField
	isDefault := address.IsDefaultField
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(addressI)).AppendWhere(_WhereFn(addressI)).Merge(
		keytitle.Update(provice.IdTitleField),
		keytitle.Update(city.IdTitleField),
		keytitle.Update(area.IdTitleField),
		phone.Update(phoneField),
		enum.Update(label),
		boolean.Update(isDefault),
		tenant.Update(address.TenatIDField),
		ownerid.Update(address.OwnerIDField),
	).AppendValidate(_DealDefault(addressI, withDefaultI))
}

func First(addressI AddressI) sqlbuilder.FirstParam {
	address := addressI.GetAddress()
	provice := address.ProvinceField
	city := address.CityField
	area := address.AreaField
	phoneField := address.ContactPhoneField
	label := address.LabelField
	isDefault := address.IsDefaultField
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		keytitle.First(provice.IdTitleField),
		keytitle.First(city.IdTitleField),
		keytitle.First(area.IdTitleField),
		phone.First(phoneField),
		enum.First(label),
		boolean.First(isDefault),
		tenant.First(address.TenatIDField),
		ownerid.First(address.OwnerIDField),
	).AppendOrder(_OrderFn(isDefault))
}

func List(addressI AddressI) sqlbuilder.ListParam {
	address := addressI.GetAddress()
	provice := address.ProvinceField
	city := address.CityField
	area := address.AreaField
	phoneField := address.ContactPhoneField
	label := address.LabelField
	isDefault := address.IsDefaultField
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		keytitle.List(provice.IdTitleField),
		keytitle.List(city.IdTitleField),
		keytitle.List(area.IdTitleField),
		phone.List(phoneField),
		enum.List(label),
		boolean.List(isDefault),
		tenant.List(address.TenatIDField),
		ownerid.List(address.OwnerIDField),
	).AppendOrder(_OrderFn(isDefault))
}

func Total(addressI AddressI) sqlbuilder.TotalParam {
	address := addressI.GetAddress()
	provice := address.ProvinceField
	city := address.CityField
	area := address.AreaField
	phoneField := address.ContactPhoneField
	label := address.LabelField
	isDefault := address.IsDefaultField
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(addressI)).Merge(
		keytitle.Total(provice.IdTitleField),
		keytitle.Total(city.IdTitleField),
		keytitle.Total(area.IdTitleField),
		phone.Total(phoneField),
		enum.Total(label),
		boolean.Total(isDefault),
		tenant.Total(address.TenatIDField),
		ownerid.Total(address.OwnerIDField),
	)
}

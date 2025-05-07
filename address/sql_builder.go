package address

import (
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit"
	"github.com/suifengpiao14/businessunit/boolean"
	"github.com/suifengpiao14/businessunit/districtcode"
	"github.com/suifengpiao14/commonlanguage"
	"github.com/suifengpiao14/sqlbuilder"
)

type AddressRule struct {
	TenatID   *sqlbuilder.Field // 业务、应用、租户等唯一标识
	OwnerID   *sqlbuilder.Field
	Label     *commonlanguage.EnumField
	MaxNumber sqlbuilder.Field // 单个业务下指定类型可配置最大条数
}

type AddressRules []AddressRule

func (rs AddressRules) GetByLabel(tenatID *sqlbuilder.Field, ownerID *sqlbuilder.Field, label *commonlanguage.EnumField) (addressRule *AddressRule, exist bool) {
	for _, r := range rs {
		if r.TenatID.IsEqual(*tenatID) && r.OwnerID.IsEqual(*ownerID) && r.Label.Field.IsEqual(*label.Field) {
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
	table        sqlbuilder.TableConfig
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
	CheckRuleI   CheckRuleI
	WithDefaultI WithDefaultI
}

type AddressFields struct {
	TenatIDField      *sqlbuilder.Field // 业务、应用、租户等唯一标识
	OwnerIDField      *sqlbuilder.Field
	LabelField        *commonlanguage.EnumField
	ContactPhoneField *sqlbuilder.Field
	ContactNameField  *sqlbuilder.Field
	AddressField      *sqlbuilder.Field
	IsDefaultField    *boolean.Boolean

	ProvinceField *districtcode.District
	CityField     *districtcode.District
	AreaField     *districtcode.District
}

func (addr AddressFields) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		addr.TenatIDField,
		addr.OwnerIDField,
		addr.LabelField.Field,
		addr.ContactPhoneField,
		addr.ContactNameField,
		addr.AddressField,
		addr.IsDefaultField.Field,
	}
	fs.Append(addr.ProvinceField.Fields()...)
	fs.Append(addr.CityField.Fields()...)
	fs.Append(addr.AreaField.Fields()...)
	return fs
}

const (
	Field_Name_isDefault = "isDefault"
)

func (addr Address) GetTable() sqlbuilder.TableConfig {
	return addr.table
}

func (addr Address) Fields() *AddressFields {
	labelField := commonlanguage.NewEnumField(addr.Label, sqlbuilder.Enums{
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
	labelField.Field.Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("label").SetTitle("标签")
	})
	addressFields := &AddressFields{
		TenatIDField: businessunit.NewTenantField(addr.TenantID),
		OwnerIDField: commonlanguage.NewOwnerID(addr.OwnerID),
		LabelField:   labelField,

		ContactPhoneField: commonlanguage.NewPhone(addr.ContactPhone).SetName("contactPhone").SetTitle("联系手机号"),
		ContactNameField: businessunit.NewNameField(addr.ContactName).SetName("contactName").SetTitle("联系人").Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.RequiredWhenInsert(true)
			f.MinBoundaryWhereInsert(1, 1)
		}).MergeSchema(sqlbuilder.Schema{
			MaxLength: 20, // 常规名称在20个字以内
			Type:      sqlbuilder.Schema_Type_string,
		}),
		AddressField: commonlanguage.NewAddress(addr.Address),
		IsDefaultField: boolean.NewBoolean(addr.IsDefault,
			sqlbuilder.Enum{
				Key:   "1",
				Title: "是",
			},
			sqlbuilder.Enum{
				Key:   "2",
				Title: "否",
			},
		).MiddlewareFn(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.SetName("isDefault").SetTitle("默认").SetFieldName(Field_Name_isDefault)
		}),
		ProvinceField: districtcode.NewDistrict(addr.ProvinceCode, addr.ProvinceName).MiddlewareFn(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			switch f.GetFieldName() {
			case districtcode.Field_Name_Code:
				f.SetName("proviceId").SetTitle("省ID")
			case districtcode.Field_Name_Name:
				f.SetName("provice").SetTitle("省").MergeSchema(sqlbuilder.Schema{
					MaxLength: 16,
				})
			}
		}),
		CityField: districtcode.NewDistrict(addr.ProvinceCode, addr.ProvinceName).MiddlewareFn(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {

			switch f.GetFieldName() {
			case districtcode.Field_Name_Code:
				f.SetName("cityId").SetTitle("城市ID")
			case districtcode.Field_Name_Name:
				f.SetName("city").SetTitle("城市").MergeSchema(sqlbuilder.Schema{
					MaxLength: 32, //海南省黎母山林场（海南黎母山省级自然保护区管理站）线上最长 25，设置32 2的倍数
				})
			}
		}),
		AreaField: districtcode.NewDistrict(addr.AreaCode, addr.AreaName).MiddlewareFn(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			switch f.GetFieldName() {
			case districtcode.Field_Name_Code:
				f.SetName("areaId").SetTitle("区ID")
			case districtcode.Field_Name_Name:
				f.SetName("area").SetTitle("区").MergeSchema(sqlbuilder.Schema{
					MaxLength: 128, //海南省黎母山林场（海南黎母山省级自然保护区管理站）线上最长 25，设置32 2的倍数
				})
			}

		}),
	}

	addressFields.TenatIDField.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(
			_ValidateRuleFn(addr.GetTable(), *addressFields, addr.CheckRuleI),
			_DealDefault(addr.GetTable(), *addressFields, addr.WithDefaultI),
		)
	})
	addressFields.TenatIDField.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(
			_DealDefault(addr.GetTable(), *addressFields, addr.WithDefaultI),
		)
	})
	bol := boolean.NewBooleanFromField(addressFields.IsDefaultField.Field)
	addressFields.IsDefaultField.Field.SetOrderFn(sqlbuilder.OrderFieldFn(bol.TrueEnum.Key, bol.FalseEnum.Key))

	return addressFields

}

func NewAddress(table sqlbuilder.TableConfig, checkRuleI CheckRuleI, withDWithDefaultI WithDefaultI) *Address {
	address := &Address{
		table:        table,
		CheckRuleI:   checkRuleI,
		WithDefaultI: withDWithDefaultI,
	}
	return address

}

// WithDefaultI 需要设置默认地址时,需要实现该接口
type WithDefaultI interface {
	CleanDefault(rawSql string) (err error)
}

type CheckRuleI interface {
	GetCount(rawSql string) (count int, err error) // 某种类型需要限制数量时,需要实现该接口,查询数据库已有的数量
}

func _ValidateRuleFn(table sqlbuilder.TableConfig, address AddressFields, checkRuleI CheckRuleI) sqlbuilder.ValueFn {
	return sqlbuilder.ValueFn{
		Layer: sqlbuilder.Value_Layer_ApiValidate,
		Fn: func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
			if checkRuleI == nil {
				return in, nil
			}
			r, ok := GetAddressRules().GetByLabel(address.TenatIDField, address.OwnerIDField, address.LabelField)
			if !ok {
				return in, nil
			}
			val, err := r.MaxNumber.GetValue(sqlbuilder.Layer_all)
			if err != nil {
				return nil, err
			}
			maxNumber := cast.ToInt(val)
			if maxNumber == 0 {
				return in, nil
			}
			rawSql, err := sqlbuilder.NewTotalBuilder(table).AppendFields(
				address.TenatIDField,
				address.OwnerIDField,
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
					address.TenatIDField.LogString(),
					address.OwnerIDField.LogString(),
					address.LabelField.Field.LogString(),
					count,
					maxNumber,
				)
				return nil, err
			}

			return in, nil
		},
	}
}

// _DealDefault 当前记录需要设置为默认记录时,清除已有的默认记录
func _DealDefault(table sqlbuilder.TableConfig, address AddressFields, withDWithDefaultI WithDefaultI) sqlbuilder.ValueFn {
	return sqlbuilder.ValueFn{
		Layer: sqlbuilder.Value_Layer_ApiFormat,
		Fn: func(val any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
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
			rawSql, err := sqlbuilder.NewUpdateBuilder(table).AppendFields(falseField.Field).AppendFields(
				address.TenatIDField,
				address.OwnerIDField,
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
		},
	}
}

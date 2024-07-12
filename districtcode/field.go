package districtcode

import (
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/businessunit/idtitle"
	"github.com/suifengpiao14/sqlbuilder"
)

/**
CREATE TABLE `t_city_info` (
  `Farea_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '城市ID',
  `Fparent_id` int(10) unsigned DEFAULT '0' COMMENT '上级城市ID',
  `Farea_name` char(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '城市名称',
  `Fall_name` char(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市全称',
  `Ffirst_letter` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '名称首字母',
  `Fcity_jianpin` char(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市简拼',
  `Fcity_qp` char(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市全拼',
  `Fcity_level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '城市级别',
  `Fcity_msg` char(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '城市信息(备注)',
  `Fcreate_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `Fupdate_time` int(10) unsigned DEFAULT '0' COMMENT '更新时间',
  `Fcity_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '城市状态信息, 1:启用,2:停用,99:删除',
  `Fprovince_id` char(12) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '省会id',
  `Fcity_id` char(12) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '城市id',
  `Fcounty_id` char(12) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '县id',
  `Fauto_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '变更时间',
  PRIMARY KEY (`Farea_id`),
  KEY `FKcityStatus` (`Fcity_status`),
  KEY `FKparentId` (`Fparent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4294967295 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地区表(主要数据来源于国家统计局并作清洗)';
**/

func NewDistrictCodeField(valueFn sqlbuilder.ValueFn) identity.IdentityField {
	field := *identity.NewIdentityField(valueFn)
	field.SetName("code").SetTitle("行政区代码").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 16, //统计局统一是使用12位，如 659008103505
		MinLength: 2,  // 只使用省时，为2位
	})

	return field
}

func NewDistrictNameField(valueFn sqlbuilder.ValueFn) sqlbuilder.Field {
	field := idtitle.NewTitleField(valueFn)
	field.SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64, //线上统计最长55个汉字，使用64符合大部分场景
		MinLength: 2,  // 只使用省时，为2个汉字
	})
	return *field
}

func NewProvinceField(codeValueFn sqlbuilder.ValueFn, nameValueFn sqlbuilder.ValueFn) *idtitle.IdTitle {
	id, title := NewDistrictCodeField(codeValueFn), NewDistrictNameField(nameValueFn)
	id.SetName("provinceCode").SetTitle("省编码")
	title.SetName("province").SetTitle("省")
	field := &idtitle.IdTitle{
		ID:    &id,
		Title: &title,
	}
	return field
}

func NewCityField(codeValueFn sqlbuilder.ValueFn, nameValueFn sqlbuilder.ValueFn) *idtitle.IdTitle {
	id, title := NewDistrictCodeField(codeValueFn), NewDistrictNameField(nameValueFn)
	id.SetName("cityCode").SetTitle("城市编码")
	title.SetName("city").SetTitle("城市")
	field := &idtitle.IdTitle{
		ID:    &id,
		Title: &title,
	}
	return field
}

func NewAreaField(codeValueFn sqlbuilder.ValueFn, nameValueFn sqlbuilder.ValueFn) *idtitle.IdTitle {
	id, title := NewDistrictCodeField(codeValueFn), NewDistrictNameField(nameValueFn)
	id.SetName("areaCode").SetTitle("城市编码")
	title.SetName("area").SetTitle("城市")
	field := &idtitle.IdTitle{
		ID:    &id,
		Title: &title,
	}
	return field
}

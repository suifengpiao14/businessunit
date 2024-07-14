package districtcode

import (
	"math"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/sqlbuilder"
)

/**
CREATE TABLE `t_city_info` (
  `Farea_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '城市ID',
  `Farea_name` char(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '城市名称',
  PRIMARY KEY (`Farea_id`),
  KEY `FKcityStatus` (`Fcity_status`),
  KEY `FKparentId` (`Fparent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4294967295 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地区表(主要数据来源于国家统计局并作清洗)';
**/

const (
	Depth_max = 5 // 省市区最大5级
)

func OptionsGetChildren(codeField *sqlbuilder.Field, nameField *sqlbuilder.Field, depth int) {
	codeField.WhereFns.AppendIfNotFirst(GetChildrenWhereFn(depth))
	nameField.WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnShield)
}

// GetChildrenWhereFn 获取子集where 函数(包含自己)depth<=0 不限制子级层级
func GetChildrenWhereFn(depth int) (whereValueFn sqlbuilder.WhereValueFn) {
	return func(dbColumnName string, in any) (value any, err error) {
		if depth <= 0 {
			depth = math.MaxInt
		}
		code := cast.ToString(in)
		dc := &_CodeLevel{}
		dc = dc.Deserialize(code) // 重新初始化化值
		return sqlbuilder.Ilike{dc.Levels.GetChildrenLikePlaceholder(depth)}, nil
	}
}

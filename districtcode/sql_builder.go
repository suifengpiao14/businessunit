package districtcode

import (
	"math"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/sqlbuilder"
)

type District struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	CodeField *sqlbuilder.Field
	NameField *sqlbuilder.Field
}

const (
	Field_Name_Code = "code"
	Field_Name_Name = "name"
)

func (d *District) Init() {
	d.CodeField = sqlbuilder.NewField(func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return d.Code, nil }).SetName("code").SetTitle("行政区代码").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 16, //统计局统一是使用12位，如 659008103505
		MinLength: 2,  // 只使用省时，为2位
	}).SetFieldName(Field_Name_Code)
	d.NameField = sqlbuilder.NewField(func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return d.Name, nil }).SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Required:  true,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 16, //统计局统一是使用12位，如 659008103505
		MinLength: 2,  // 只使用省时，为2位
	}).SetFieldName(Field_Name_Name)
}

func (d District) Fields() (fs sqlbuilder.Fields) {
	fs = sqlbuilder.Fields{
		d.CodeField,
		d.NameField,
	}

	return fs
}

func (p *District) MiddlewareFn(initFns ...sqlbuilder.ApplyFn) *District {
	p.CodeField.Apply(initFns...)
	return p
}

func NewDistrict(code string, name string) *District {
	d := &District{
		Code: code,
		Name: name,
	}
	d.Init()
	return d
}

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
	codeField.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil, GetChildrenWhereFn(depth))
	nameField.WhereFns.Append(sqlbuilder.ValueFnShield)
}

// GetChildrenWhereFn 获取子集where 函数(包含自己)depth<=0 不限制子级层级
func GetChildrenWhereFn(depth int) (whereValueFn sqlbuilder.ValueFn) {
	return sqlbuilder.ValueFn{
		Fn: func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (value any, err error) {
			if depth <= 0 {
				depth = math.MaxInt
			}
			code := cast.ToString(in)
			dc := &_CodeLevel{}
			dc = dc.Deserialize(code) // 重新初始化化值
			return sqlbuilder.Ilike{dc.Levels.GetChildrenLikePlaceholder(depth)}, nil
		},
		Layer: sqlbuilder.Value_Layer_SetValue,
	}
}

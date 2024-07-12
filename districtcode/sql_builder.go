package districtcode

import (
	"math"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/idtitle"
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

type District struct {
	*DistrictField
}

func (address District) CustomFields() (fileds sqlbuilder.Fields) {
	fileds = make(sqlbuilder.Fields, 0)
	return fileds
}

type DistrictI interface {
	GetDistrict() *District
	sqlbuilder.TableI
}

func _DataFn(districtI DistrictI) sqlbuilder.DataFn {
	return districtI.GetDistrict().CustomFields().Data
}

func _WhereFn(districtI DistrictI) sqlbuilder.WhereFn {
	return districtI.GetDistrict().CustomFields().Where
}

func Insert(districtI DistrictI) sqlbuilder.InsertParam {
	district := districtI.GetDistrict()
	idtitleField := district.DistrictField.IdTitleField
	return sqlbuilder.NewInsertBuilder(nil).AppendData(_DataFn(districtI)).Merge(
		idtitle.Insert(idtitleField),
	)
}

func Update(districtI DistrictI) sqlbuilder.UpdateParam {
	district := districtI.GetDistrict()
	idtitleField := district.DistrictField.IdTitleField
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(_DataFn(districtI)).AppendWhere(_WhereFn(districtI)).Merge(
		idtitle.Update(idtitleField),
	)
}

func First(districtI DistrictI) sqlbuilder.FirstParam {
	district := districtI.GetDistrict()
	idtitleField := district.DistrictField.IdTitleField
	return sqlbuilder.NewFirstBuilder(nil).AppendWhere(_WhereFn(districtI)).Merge(
		idtitle.First(idtitleField),
	)
}

func List(districtI DistrictI) sqlbuilder.ListParam {
	district := districtI.GetDistrict()
	idtitleField := district.DistrictField.IdTitleField
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(districtI)).Merge(
		idtitle.List(idtitleField),
	)
}

func Total(districtI DistrictI) sqlbuilder.TotalParam {
	district := districtI.GetDistrict()
	idtitleField := district.DistrictField.IdTitleField
	return sqlbuilder.NewTotalBuilder(nil).AppendWhere(_WhereFn(districtI)).Merge(
		idtitle.Total(idtitleField),
	)
}

// GetChildrenWhereFn 获取子集where 函数(包含自己)depth<=0 不限制子级层级
func GetChildrenWhereFn(depth int) (whereValueFn sqlbuilder.ValueFn) {
	return func(in any) (value any, err error) {
		if depth <= 0 {
			depth = math.MaxInt
		}
		code := cast.ToString(in)
		dc := &_CodeLevel{}
		dc = dc.Deserialize(code) // 重新初始化化值
		return sqlbuilder.Ilike{dc.Levels.GetChildrenLikePlaceholder(depth)}, nil
	}
}

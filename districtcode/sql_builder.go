package districtcode

import (
	"fmt"
	"math"
	"sort"
	"strings"

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

func GetByParentCode(districtI DistrictI, depth int) sqlbuilder.ListParam {
	districtI.GetDistrict().DistrictField.ID.Field.WhereFns.AppendIfNotFirst(func(in any) (any, error) {
		parentCode := cast.ToString(in)
		parentCode = strings.TrimSuffix(parentCode, "0") // 删除结尾的0, 省留下留2位，市留下4位，区/县留下6位 街道留下9位置 小区留下12位置
		return sqlbuilder.Ilike{parentCode, "%"}, nil
	})
	return sqlbuilder.NewListBuilder(nil).AppendWhere(_WhereFn(districtI))
}

// Level 等级
type Level struct {
	Code  string `json:"code"`
	Level int    `json:"level"`
}

// LikeChart like 语句中匹配本级字符
func (l Level) LikeChart() string {
	s := strings.Repeat("_", len(l.Code))
	return s
}

type Levels [5]Level

func (a Levels) Len() int           { return len(a) }
func (a Levels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Levels) Less(i, j int) bool { return a[i].Level < a[j].Level }

func (ls Levels) IsChildrenWithSelf(target Levels, depth int) (ok bool) {
	sort.Sort(ls)
	sort.Sort(target)
	for i := 0; i < 5; i++ {
		clLevel := ls[i]
		targetLevel := target[i]
		if clLevel == targetLevel {
			continue // 同级相同
		}
		//不同
		if depth < 1 {
			return false // 超过指定级别后不计算为子类
		}
		depth-- //当前级别容许不同
		clLevelInt := cast.ToInt(clLevel.Code)
		targetLevelInt := cast.ToInt(targetLevel.Code)
		if targetLevelInt < clLevelInt {
			return false // 父类值小于子级
		}
	}
	return true
}

func (ls Levels) GetChildrenLikePlaceholder(depth int) (placeholder string) {
	sort.Sort(ls)
	arr := make([]string, 0)
	for _, l := range ls {
		codeInt := cast.ToInt(l.Code)
		if codeInt == 0 && depth > 0 { // 子集替换
			depth--
			arr = append(arr, l.LikeChart())
			continue
		}
		arr = append(arr, l.Code)
	}

	placeholder = strings.Join(arr, "")
	return placeholder

}

type CodeLevel struct {
	Levels Levels
	ref    RecordI
}

// Level  获取区域码级别0-全国,1-省,2-市,3-县,4-乡/镇,5-村/居委会/小区
func (cl CodeLevel) Level() int {
	sort.Sort(cl.Levels)
	for i, l := range cl.Levels {
		codeInt := cast.ToInt(l.Code)
		if codeInt == 0 {
			return i
		}
	}
	return 0
}

func (cl CodeLevel) IsSame(target CodeLevel) (ok bool) {
	return strings.EqualFold(cl.Serialize(), target.Serialize())
}

func (cl CodeLevel) IsChildrenWithSelf(target CodeLevel, depth int) (ok bool) {
	return cl.Levels.IsChildrenWithSelf(target.Levels, depth)
}

func (cl CodeLevel) Serialize() (s string) {
	arr := make([]string, 0)
	sort.Sort(cl.Levels)
	for _, l := range cl.Levels {
		arr = append(arr, l.Code)
	}
	s = strings.Join(arr, "")
	return s
}

const (
	Level_Level_country  = iota // 国家级
	Level_Level_province        // 省级
	Level_Level_city            // 市级
	Level_Level_county          // 区/县级
	Level_Level_town            // 乡/镇/街道级别
	Level_Level_village         // 村/小区级别
)

// Deserialize 将code 按照 省(2位)市(2位)区/县(2位)乡/镇/街道(3位)村/小区/居委会(3位)格式返回
func (dc *CodeLevel) Deserialize(code string) *CodeLevel {
	code = fmt.Sprintf("%s%s", code, strings.Repeat("0", 12)) // 先补12个0
	code = code[:12]                                          // 截取前12位
	province, city, county, town, village := code[:2], code[2:4], code[4:6], code[6:9], code[9:12]
	*dc = CodeLevel{
		Levels: Levels{
			{Code: province, Level: Level_Level_province},
			{Code: city, Level: Level_Level_city},
			{Code: county, Level: Level_Level_county},
			{Code: town, Level: Level_Level_town},
			{Code: village, Level: Level_Level_village},
		},
	}

	return dc
}

// GetChildrenWhereFn 获取子集where 函数(包含自己)depth<=0 不限制子级层级
func GetChildrenWhereFn(depth int) (whereValueFn sqlbuilder.ValueFn) {
	return func(in any) (value any, err error) {
		if depth <= 0 {
			depth = math.MaxInt
		}
		code := cast.ToString(in)
		dc := &CodeLevel{}
		dc = dc.Deserialize(code) // 重新初始化化值
		return sqlbuilder.Ilike{dc.Levels.GetChildrenLikePlaceholder(depth)}, nil
	}
}

type CodeLevels []*CodeLevel

func (a CodeLevels) Len() int           { return len(a) }
func (a CodeLevels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CodeLevels) Less(i, j int) bool { return a[i].Level() < a[j].Level() }

func (cls CodeLevels) getAllRef() (refs []RecordI) {
	refs = make([]RecordI, 0)
	for _, cl := range cls {
		refs = append(refs, cl.ref)
	}
	return refs
}

// GetChildren 获取子级
func (cls CodeLevels) GetChildren(parent CodeLevel) (children CodeLevels) {
	children = make(CodeLevels, 0)
	for _, cl := range cls {
		if parent.IsSame(*cl) {
			continue
		}
		if parent.IsChildrenWithSelf(*cl, 1) { // 相邻下一层
			children = append(children, cl)
		}
	}
	return children
}

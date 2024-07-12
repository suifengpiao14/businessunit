package districtcode

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type RecordI interface {
	GetCode() string
	AddChildren(children ...RecordI)
}

// Tree 行列式改成Tree 格式
func Tree(records []RecordI) (trees []RecordI) {
	trees = make([]RecordI, 0)
	cls := make(_CodeLevels, 0)
	for _, record := range records {
		cl := &_CodeLevel{}
		cl = cl.Deserialize(record.GetCode())
		cl.ref = record
		cls = append(cls, cl)
	}
	for _, cl := range cls {
		children := cls.GetChildren(*cl)
		cl.ref.AddChildren(children.getAllRef()...)
	}
	sort.Sort(cls)
	topLevel := 0
	for i, cl := range cls {
		if i == 0 {
			topLevel = cl.Level()
		}
		if topLevel == cl.Level() {
			trees = append(trees, cl.ref)
			continue
		}
		return trees

	}
	return trees
}

// _Level 等级
type _Level struct {
	Code  string `json:"code"`
	Level int    `json:"level"`
}

// LikeChart like 语句中匹配本级字符
func (l _Level) LikeChart() string {
	s := strings.Repeat("_", len(l.Code))
	return s
}

type _Levels [5]_Level

func (a _Levels) Len() int           { return len(a) }
func (a _Levels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a _Levels) Less(i, j int) bool { return a[i].Level < a[j].Level }

func (ls _Levels) IsChildrenWithSelf(target _Levels, depth int) (ok bool) {
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

func (ls _Levels) GetChildrenLikePlaceholder(depth int) (placeholder string) {
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

type _CodeLevel struct {
	Levels _Levels
	ref    RecordI
}

// Level  获取区域码级别0-全国,1-省,2-市,3-县,4-乡/镇,5-村/居委会/小区
func (cl _CodeLevel) Level() int {
	sort.Sort(cl.Levels)
	for i, l := range cl.Levels {
		codeInt := cast.ToInt(l.Code)
		if codeInt == 0 {
			return i
		}
	}
	return 0
}

func (cl _CodeLevel) IsSame(target _CodeLevel) (ok bool) {
	return strings.EqualFold(cl.Serialize(), target.Serialize())
}

func (cl _CodeLevel) IsChildrenWithSelf(target _CodeLevel, depth int) (ok bool) {
	return cl.Levels.IsChildrenWithSelf(target.Levels, depth)
}

func (cl _CodeLevel) Serialize() (s string) {
	arr := make([]string, 0)
	sort.Sort(cl.Levels)
	for _, l := range cl.Levels {
		arr = append(arr, l.Code)
	}
	s = strings.Join(arr, "")
	return s
}

const (
	_Level_Level_country  = iota // 国家级
	_Level_Level_province        // 省级
	_Level_Level_city            // 市级
	_Level_Level_county          // 区/县级
	_Level_Level_town            // 乡/镇/街道级别
	_Level_Level_village         // 村/小区级别
)

// Deserialize 将code 按照 省(2位)市(2位)区/县(2位)乡/镇/街道(3位)村/小区/居委会(3位)格式返回
func (dc *_CodeLevel) Deserialize(code string) *_CodeLevel {
	code = fmt.Sprintf("%s%s", code, strings.Repeat("0", 12)) // 先补12个0
	code = code[:12]                                          // 截取前12位
	province, city, county, town, village := code[:2], code[2:4], code[4:6], code[6:9], code[9:12]
	*dc = _CodeLevel{
		Levels: _Levels{
			{Code: province, Level: _Level_Level_province},
			{Code: city, Level: _Level_Level_city},
			{Code: county, Level: _Level_Level_county},
			{Code: town, Level: _Level_Level_town},
			{Code: village, Level: _Level_Level_village},
		},
	}

	return dc
}

type _CodeLevels []*_CodeLevel

func (a _CodeLevels) Len() int           { return len(a) }
func (a _CodeLevels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a _CodeLevels) Less(i, j int) bool { return a[i].Level() < a[j].Level() }

func (cls _CodeLevels) getAllRef() (refs []RecordI) {
	refs = make([]RecordI, 0)
	for _, cl := range cls {
		refs = append(refs, cl.ref)
	}
	return refs
}

// GetChildren 获取子级
func (cls _CodeLevels) GetChildren(parent _CodeLevel) (children _CodeLevels) {
	children = make(_CodeLevels, 0)
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

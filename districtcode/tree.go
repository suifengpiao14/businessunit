package districtcode

import "sort"

type RecordI interface {
	GetCode() string
	AddChildren(children ...RecordI)
}

// Tree 行列式改成Tree 格式
func Tree(records []RecordI) (trees []RecordI) {
	trees = make([]RecordI, 0)
	cls := make(CodeLevels, 0)
	for _, record := range records {
		cl := &CodeLevel{}
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

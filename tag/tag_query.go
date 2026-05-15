package tag

import (
	"gitlab.huishoubao.com/gopackage/sqlbuilder"
)

type QueryAll struct {
	Dimension *sqlbuilder.Field
	Table     string
	Builder   sqlbuilder.Builder
}

func (q QueryAll) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{q.Dimension}
}

func (q QueryAll) Query(result any) (err error) {
	return q.Builder.ListParam(q.Fields()).List(result)
}

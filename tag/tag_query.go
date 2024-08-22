package tag

import (
	"github.com/suifengpiao14/sqlbuilder"
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
	return q.Builder.List(result, q.Fields()...)
}

package tag

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type QueryAll struct {
	Dimension    *sqlbuilder.Field
	Table        string
	QueryHandler sqlbuilder.QueryHandler
}

func (q QueryAll) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{q.Dimension}
}

func (q QueryAll) Param() (builder sqlbuilder.ListParam) {
	return sqlbuilder.NewListBuilder(q.Table).AppendFields(q.Fields()...)
}

func (q QueryAll) Query(result any) (err error) {
	param := q.Param()
	return param.Query(result, q.QueryHandler)
}

package column

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

func init() {
	IdentityColumnDefault.UpdateWhere = func() (expressions []goqu.Expression, err error) {
		val, err := IdentityColumnDefault.UpdateData()
		if err != nil {
			return nil, err
		}
		expressions = sqlbuilder.ConcatExpression(goqu.Ex{IdentityColumnDefault.Name: val})
		return

	}
}

var IdentityColumnDefault = sqlbuilder.Column{
	Name: "id",
	InsertData: func() (value any, err error) {
		id := xid.New().String()
		return id, nil
	},
}

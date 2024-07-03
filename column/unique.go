package column

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/suifengpiao14/sqlbuilder"
)

type UniqueColumn sqlbuilder.Multicolumn

func (uc UniqueColumn) UpdateWhere() (expressions []goqu.Expression, err error) {
	expressions = make([]goqu.Expression, 0)
	for _, c := range uc.Columns {
		val, err := c.UpdateData()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, goqu.C(c.Name).Eq(val))
	}
	return expressions, nil

}

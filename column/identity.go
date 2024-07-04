package column

import (
	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

type IdentityColumn sqlbuilder.Column

func (c IdentityColumn) InsertValue() (value any, err error) {
	if c.InsertValueFn != nil {
		return c.InsertValueFn()
	}
	id := xid.New().String()
	return id, nil
}

func (c IdentityColumn) UpdateValue() (value any, err error) {
	return // 重写，确保不可修改
}

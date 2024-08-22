package tag

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type CmdAdd struct {
	Dimension    *sqlbuilder.Field // 该列可能为空
	Tag          sqlbuilder.Field
	Table        string
	ExtendFields func(cmdAdd *CmdAdd) sqlbuilder.Fields
	Builder      sqlbuilder.Builder
}

func (q CmdAdd) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{}

	if q.Dimension != nil {
		q.Dimension.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.Append(sqlbuilder.ValueFnForward)
		})
	}
	q.Tag.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnForward)
	})
	fs.Append(&q.Tag, q.Dimension)
	if q.ExtendFields != nil {
		fs.Append(q.ExtendFields(&q)...)
	}

	return sqlbuilder.Fields{&q.Tag, q.Dimension}
}

func (cmd CmdAdd) Exec() (err error) {

	exists, err := cmd.Builder.Exists(cmd.Fields()...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = cmd.Builder.Insert(cmd.Fields()...)
	return err

}

package tag

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type CmdAdd struct {
	Classify     *sqlbuilder.Field // 该列可能为空
	Tag          sqlbuilder.Field
	Table        string
	ExtendFields func(cmdAdd *CmdAdd) sqlbuilder.Fields
	QueryHandler sqlbuilder.QueryHandler // 查询是否存在
	ExecHandler  sqlbuilder.ExecHandler  //执行插入
}

func (q CmdAdd) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{}

	if q.Classify != nil {
		q.Classify.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.Append(sqlbuilder.ValueFnForward)
		})
	}
	q.Tag.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnForward)
	})
	fs.Append(&q.Tag, q.Classify)
	if q.ExtendFields != nil {
		fs.Append(q.ExtendFields(&q)...)
	}

	return sqlbuilder.Fields{&q.Tag, q.Classify}
}

func (action CmdAdd) ExistsBuilder() (builder sqlbuilder.ExistsParam) {
	return sqlbuilder.NewExistsBuilder(action.Table).AppendFields(action.Fields()...)
}

func (action CmdAdd) InsertBuilder() (builder sqlbuilder.InsertParam) {
	return sqlbuilder.NewInsertBuilder(action.Table).AppendFields(action.Fields()...)
}
func (action CmdAdd) Exists() (exists bool, err error) {
	return action.ExistsBuilder().Exists(action.QueryHandler)
}
func (action CmdAdd) Exec() (err error) {
	exists, err := action.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = action.InsertBuilder().Exec(action.ExecHandler)
	return err

}

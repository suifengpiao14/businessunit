package tag

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type CmdAdd struct {
	Dimension *sqlbuilder.Field // 该列可能为空
	Tag       *sqlbuilder.Field
	Table     string
	Builder   sqlbuilder.Builder
}

func (q CmdAdd) Fields() sqlbuilder.Fields {
	dimensionField := sqlbuilder.NewStringField("", "dimension", "维度", 0).SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnForward)
	})

	tagField := sqlbuilder.NewStringField("", "tag", "标签", 0).SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnForward)
	})

	fs := sqlbuilder.Fields{}
	if q.Dimension != nil {
		q.Dimension.Combine(dimensionField)
		fs = append(fs, q.Dimension)
	}
	if q.Tag != nil {
		q.Tag.Combine(tagField)
		fs = append(fs, q.Tag)
	}

	return fs
}

func (cmd CmdAdd) Exec(fields ...*sqlbuilder.Field) (err error) {

	exists, err := cmd.Builder.Exists(fields...)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = cmd.Builder.Insert(fields...)
	return err

}

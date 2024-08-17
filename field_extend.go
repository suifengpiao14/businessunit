package businessunit

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

func NewNickname(nickname string) *sqlbuilder.Field {
	f := NewNameField(nickname).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetName("nickname").SetTitle("昵称")
	})
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(
			sqlbuilder.ValueFnEmpty2Nil,
			sqlbuilder.ValueFnWhereLike,
		)
	})
	return f
}

type KeyTitleFields[keyT int | string] struct {
	Key   sqlbuilder.FieldFn[keyT]
	Title sqlbuilder.FieldFn[string]
}

func (KeyTitleFields[T]) Fields() KeyTitleFields[T] {

	return KeyTitleFields[T]{
		Key:   NewKeyField[T],
		Title: NewTitleField,
	}
}

type IdNameFields[idT int | string] struct {
	IdField   sqlbuilder.FieldFn[idT]
	NameField sqlbuilder.FieldFn[string]
}

func (IdNameFields[T]) Builder() IdNameFields[T] {
	return IdNameFields[T]{
		IdField: func(value T) *sqlbuilder.Field {
			return NewIdentifierField(value)
		},
		NameField: func(value string) *sqlbuilder.Field {
			return NewNameField(value)
		},
	}
}

type CUTimeFields struct {
	CreatedAt sqlbuilder.FieldFn[string]
	UpdatedAt sqlbuilder.FieldFn[string]
}

func (CUTimeFields) Builder() CUTimeFields {
	return CUTimeFields{
		CreatedAt: func(value string) *sqlbuilder.Field {
			return NewCreatedAtField()
		},
		UpdatedAt: func(value string) *sqlbuilder.Field {
			return NewUpdatedAtField()
		},
	}
}

type CUDTimeFields struct {
	CreatedAt sqlbuilder.FieldFn[string]
	UpdatedAt sqlbuilder.FieldFn[string]
	DeletedAt sqlbuilder.FieldFn[string]
}

func (CUDTimeFields) Builder() CUDTimeFields {
	cuTime := new(CUTimeFields).Builder()
	return CUDTimeFields{
		CreatedAt: cuTime.CreatedAt,
		UpdatedAt: cuTime.UpdatedAt,
		DeletedAt: func(value string) *sqlbuilder.Field {
			return NewDeletedAtField()
		},
	}
}

// UniqueFieldMiddlewarSceneInsert 单列唯一索引键,新增场景中间件
func UniqueFieldMiddlewarSceneInsert(table string, existsFn func(sql string) (exists bool, err error)) sqlbuilder.MiddlewareFn {
	return func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		sceneFnName := "checkexists"
		sceneFn := sqlbuilder.SceneFn{
			Name:  sceneFnName,
			Scene: sqlbuilder.SCENE_SQL_INSERT,
			Fn: func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
				f1 := f.Copy()               //复制不影响外部,在内部copy 是运行时 copy,确保 builder阶段的设置都能考呗到
				f1.SceneFnRmove(sceneFnName) // 避免死循环
				f1.WhereFns.Append(sqlbuilder.ValueFnForward)
				f.ValueFns.Append(func(inputValue any) (any, error) {
					totalParam := sqlbuilder.NewTotalBuilder(table).AppendFields(f1)
					sql, err := totalParam.ToSQL()
					if err != nil {
						return nil, err
					}
					exists, err := existsFn(sql)
					if err != nil {
						return nil, err
					}
					if exists {
						err = errors.Errorf("unique column %s value %s exists", f1.DBName(), inputValue)
						return nil, err
					}
					return inputValue, nil
				})

			},
		}
		f.SceneFn(sceneFn)

	}
}

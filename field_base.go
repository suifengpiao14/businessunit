package businessunit

import (
	"time"

	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

func NewNameField(name string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return name, nil }).SetName("name").SetTitle("名称").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
	})

	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.MergeSchema(sqlbuilder.Schema{Minimum: 1})
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnWhereLike)
	})
	return f
}

func NewTextField(text string, maxLength int) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return text, nil }).SetName("text").SetTitle("文本").MergeSchema(sqlbuilder.Schema{
		Type: sqlbuilder.Schema_Type_string,
	})
	if maxLength > 0 {
		f.MergeSchema(sqlbuilder.Schema{MaxLength: maxLength})
	}
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	return f
}

func NewProfile(profile string) (f *sqlbuilder.Field) {
	f = NewTextField(profile, 300).SetName("profile").SetTitle("简介")
	return f
}

var NewIntField = sqlbuilder.NewIntField
var NewStringField = sqlbuilder.NewStringField

func NewUpdatedAtField() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(time.DateTime), nil
	})
	f.SetName("updatedAt").SetTitle("更新时间").SetTag(sqlbuilder.Tag_updatedAt)
	return f
}

func NewPageIndexField(pageIndex string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(pageIndex).SetName("pageIndex").SetTitle("页码").Apply(sqlbuilder.ApplyFnValueEmpty2Nil).SetTag(sqlbuilder.Field_tag_pageIndex).SetType(sqlbuilder.Schema_Type_int).SceneInit(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Reset() // 屏蔽where 条件
	})
	return f
}

func NewPageSizeField(pageSize string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(pageSize).SetName("pageSize").SetTitle("每页数量").Apply(sqlbuilder.ApplyFnValueEmpty2Nil).SetTag(sqlbuilder.Field_tag_pageSize).SetType(sqlbuilder.Schema_Type_int).SceneInit(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Reset() // 屏蔽where 条件
	})
	return f
}

func NewCreatedAtField() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(time.DateTime), nil
	}).SetName("createdAt").SetTitle("创建时间").SetTag(sqlbuilder.Tag_createdAt)
	f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield) // 更新时屏蔽
	})
	return f
}

func NewAutoId(autoId uint) (field *sqlbuilder.Field) {
	field = sqlbuilder.NewField(func(in any) (any, error) { return autoId, nil })
	field.SetName("id").SetTitle("ID").MergeSchema(sqlbuilder.Schema{
		Type:          sqlbuilder.Schema_Type_int,
		Maximum:       sqlbuilder.Int_maximum_bigint,
		MaxLength:     64,
		Primary:       true,
		AutoIncrement: true,
	})

	field.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield)
	})
	field.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // id 不能更新
		f.WhereFns.Append(sqlbuilder.ValueFnFormatArray)
		f.SetRequired(true)
		f.MergeSchema(sqlbuilder.Schema{
			Minimum: 1,
		})
	})

	field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil, sqlbuilder.ValueFnFormatArray)
		if f.Schema.Required {
			f.MergeSchema(sqlbuilder.Schema{
				Minimum: 1,
			})
		}
	})
	return field
}

func NewTenantField[T int | string](tenant T) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return tenant, nil }).SetName("ternatId").SetTitle("租户ID")
	f.MergeSchema(sqlbuilder.Schema{
		Required:  true,
		MinLength: 1,
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Maximum:   sqlbuilder.UnsinedInt_maximum_bigint,
		Minimum:   1,
	})
	f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // 不可更新
	})
	f.WhereFns.Append(sqlbuilder.ValueFnForward) // update,select 都必须为条件
	return f
}

const (
	Tag_DeletedStatusField = "DeletedStatusField"
)

// NewDeletedStatusField 使用特定状态标记删除
func NewDeletedStatusField[T int | string](deletedStatus T) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return deletedStatus, nil
	}).SetName("status").SetTitle("状态").SetTag(Tag_DeletedStatusField) // 设置特殊标记,方便使用时获取列特殊处理
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield)
	})
	f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true)
	})
	//设置删除场景
	f.SceneFn(sqlbuilder.SceneFn{
		Scene: sqlbuilder.SCENE_SQL_DELETE,
		Fn: func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.ValueFns.ResetSetValueFn(func(in any) (any, error) {
				return deletedStatus, nil
			})
		},
	})
	f.WhereFns.Append(sqlbuilder.ValueFn{
		Layer: sqlbuilder.Value_Layer_DBFormat,
		Fn: func(in any) (any, error) {
			return sqlbuilder.Neq{Value: in}, nil
		},
	})

	return f
}

func NewKeyFieldField[T int | uint | int64 | string](value T) *sqlbuilder.Field {
	f := sqlbuilder.NewField(func(in any) (any, error) {
		return value, nil
	}).SetName("key").SetTitle("键")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	f.MinBoundaryWhereInsert(1, 1)
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

func NewUuidField[T int | int64 | string](value T) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return value, nil }).SetName("uuid").SetTitle("UUID")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		MinLength: 1,
	})
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetRequired(true)
		f.ValueFns.Append(sqlbuilder.ValueFn{
			Layer: sqlbuilder.Value_Layer_DBFormat,
			Fn: func(in any) (any, error) {
				return xid.New().String(), nil
			},
		})
		f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.ShieldUpdate(true) // uuid 不能更新
			f.WhereFns.Append(sqlbuilder.ValueFnForward)
		})

		f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.Append(sqlbuilder.ValueFnForward)
		})
	})

	return f
}

func NewKeyField[T int | int64 | string](value T) *sqlbuilder.Field {
	f := sqlbuilder.NewField(func(in any) (any, error) { return value, nil }).SetName("key").SetTitle("键")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		Minimum:   1,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)

	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

func OptionForeignkey(f *sqlbuilder.Field, redundantFields ...sqlbuilder.Field) {
	if len(redundantFields) > 0 {
		return
	}
	f.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFn{
			Layer: sqlbuilder.Value_Layer_ApiFormat,
			Fn: func(in any) (any, error) {
				val, err := f.GetValue()
				if err != nil {
					return nil, err
				}
				m := map[string]any{}
				for _, redundantField := range redundantFields {
					redundantField.ValueFns.Append(sqlbuilder.ValueFn{
						Layer: sqlbuilder.Value_Layer_SetValue,
						Fn:    func(in any) (any, error) { return val, nil },
					})
					redundantFiledValue, err := redundantField.GetValue()
					if err != nil {
						return nil, err
					}
					if !sqlbuilder.IsNil(redundantFiledValue) {
						m[redundantField.DBName()] = redundantFiledValue
					}
				}
				return m, nil
			},
		})
	})

}

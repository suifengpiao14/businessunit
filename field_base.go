package businessunit

import (
	"time"

	"github.com/rs/xid"
	"github.com/suifengpiao14/sqlbuilder"
)

func NewPhoneField(phone string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return phone, nil })
	f.SetName("phone").SetTitle("手机号").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 15,
		RegExp:    `^1[3-9]\d{9}$`, // 中国大陆手机号正则表达式
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

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

func NewProfileField(profile string) (f *sqlbuilder.Field) {
	f = NewTextField(profile, 300).SetName("profile").SetTitle("简介")
	return f
}

var NewIntField = sqlbuilder.NewIntField
var NewStringField = sqlbuilder.NewStringField

func NewAddressField(address string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return address, nil }).SetName("address").SetTitle("地址").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 128, // 线上统计最大55个字符，设置128 应该适合大部分场景大小
	})
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	return f
}

func NewHeightField(height int) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return height, nil }).SetName("height").SetTitle("高").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_int,
		MaxLength: 10000, //日常物体、人、动物高不操过1万m/cm
	})
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	return f
}

func NewGenderField[T int | string](val T, man T, woman T) *EnumField {
	genderField := NewEnumField(val, sqlbuilder.Enums{
		sqlbuilder.Enum{
			Key:   man,
			Title: "男",
		},
		sqlbuilder.Enum{
			Key:   woman,
			Title: "女",
		},
	})
	genderField.Field.SetName("gender").SetTitle("性别").Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.Schema.Enums.Append(sqlbuilder.Enum{
				Key:   "",
				Title: "全部",
			})
			f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		})
	})
	return genderField
}

func NewBooleanField[T int | string](val T, enumTrue T, enumFalse T) *EnumField {
	genderField := NewEnumField(val, sqlbuilder.Enums{
		sqlbuilder.Enum{
			Key:   enumTrue,
			Title: "真",
		},
		sqlbuilder.Enum{
			Key:   enumFalse,
			Title: "假",
		},
	})
	genderField.Field.SetName("bool").SetTitle("真假").Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		})
	})
	return genderField
}

func NewEmailField(email string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return email, nil }).SetName("email").SetTitle("邮箱")
	f.MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 32,
		MinLength: 5,
		RegExp:    `([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}`, // 邮箱验证表达式
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	})
	return f
}

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

func NewAutoIdField(autoId uint) (field *sqlbuilder.Field) {
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

// NewDeletedAtField 通过删除时间列标记删除
func NewDeletedAtField() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField("").SetName("deleted_at").SetTitle("删除时间")
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
				return time.Now().Local().Format(time.DateTime), nil
			})
		},
	})
	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.ResetSetValueFn(func(inputValue any) (any, error) {
			return "", nil
		})
		f.WhereFns.Append(sqlbuilder.ValueFnForward)
	})

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
			return sqlbuilder.Neq(in), nil
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

func NewOwnerID[T int | string | int64](value T) *sqlbuilder.Field {
	field := sqlbuilder.NewField(func(in any) (any, error) { return value, nil }).SetName("ownerId").SetTitle("所有者").MergeSchema(sqlbuilder.Schema{
		Comment:      "对象标识,缺失时记录无意义",
		Type:         sqlbuilder.Schema_Type_string,
		MaxLength:    64,
		MinLength:    1,
		Minimum:      1,
		ShieldUpdate: true, // 所有者不可跟新
	})
	field.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.SetRequired(true)
	})
	field.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	return field
}

func NewIdentifierField(value any) *sqlbuilder.Field {
	f := sqlbuilder.NewField(func(in any) (any, error) { return value, nil }).SetName("identity").SetTitle("标识")
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

func NewTitleField(value string) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return value, nil })
	f.SetName("title").SetTitle("标题").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
	}).ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)

	f.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnWhereLike)
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

type EnumField struct {
	Enums sqlbuilder.Enums
	Field *sqlbuilder.Field
}

func (b *EnumField) MiddlewareFn(initFns ...sqlbuilder.ApplyFn) *EnumField {
	b.Field.Apply(initFns...)
	return b
}

func NewEnumField(value any, enums sqlbuilder.Enums) *EnumField {
	e := &EnumField{
		Enums: enums,
	}
	e.Field = sqlbuilder.NewField(func(in any) (any, error) { return value, nil }).SetName("enum_column").SetTag("枚举列")
	e.Field.AppendEnum(enums...)
	return e
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

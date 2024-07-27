package businessunit

import (
	"time"

	"github.com/suifengpiao14/businessunit/enum"
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

func NewIntField(name string, maximum uint) (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) { return name, nil }).SetName("number").SetTitle("数字").MergeSchema(sqlbuilder.Schema{
		Type: sqlbuilder.Schema_Type_int,
	})
	if maximum > 0 {
		f.MergeSchema(sqlbuilder.Schema{Maximum: maximum})
	}
	f.ValueFns.Append(sqlbuilder.ValueFnEmpty2Nil)
	return f
}

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

var NewEnum = enum.NewEnum

func NewGenderField[T int | string](val T, man T, woman T) *enum.Enum {
	genderField := enum.NewEnum(val, sqlbuilder.Enums{
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

var Time_format = sqlbuilder.Time_format

func NewUpdatedAtField() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(Time_format), nil
	})
	f.SetName("updated_at").SetTitle("更新时间")
	return f
}

func NewCreatedAt() (f *sqlbuilder.Field) {
	f = sqlbuilder.NewField(func(in any) (any, error) {
		return time.Now().Local().Format(Time_format), nil
	}).SetName("created_at").SetTitle("创建时间")
	f.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield) // 更新时屏蔽
	})
	return f
}

func NewAutoIdField(autoId uint) (field *sqlbuilder.Field) {
	field = sqlbuilder.NewField(func(in any) (any, error) { return autoId, nil })
	field.SetName("id").SetTitle("ID").MergeSchema(sqlbuilder.Schema{
		Type:      sqlbuilder.Schema_Type_int,
		MaxLength: 64,
		Primary:   true,
	})

	field.SceneInsert(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ValueFns.Append(sqlbuilder.ValueFnShield)
	})
	field.SceneUpdate(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.ShieldUpdate(true) // id 不能更新
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		if f.Schema.Required {
			f.MergeSchema(sqlbuilder.Schema{
				Minimum: 1,
			})
		}
	})

	field.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		f.WhereFns.Append(sqlbuilder.ValueFnEmpty2Nil)
		if f.Schema.Required {
			f.MergeSchema(sqlbuilder.Schema{
				Minimum: 1,
			})
		}
	})
	return field
}

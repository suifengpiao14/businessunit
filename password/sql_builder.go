package password

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/sqlbuilder"
)

type Password struct {
	UserId        string `json:"userId"`
	Password      string `json:"password"`
	UserIdField   *sqlbuilder.Field
	PasswordField *sqlbuilder.Field
}

func (p Password) Init() {
	p.PasswordField = sqlbuilder.NewField(func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) { return p.Password, nil }).SetName("password").SetTitle("密码")
	p.PasswordField.ValueFns.Append(sqlbuilder.ValueFn{
		Layer: sqlbuilder.Value_Layer_DBFormat,
		Fn: func(in any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
			password := cast.ToString(in)
			if password == "" {
				err := errors.Errorf("password request string")
				return nil, err
			}
			value := EncodingFn(password) // 对密码字段加密转换
			return value, nil
		},
	})
	p.PasswordField.MergeSchema(sqlbuilder.Schema{
		Title:     "密码",
		Required:  true,
		Comment:   "对象标识,缺失时记录无意义",
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		MinLength: 1,
	})
	p.PasswordField.WhereFns.Append(sqlbuilder.ValueFnShield)

}

func (p Password) Fields() (fs sqlbuilder.Fields) {
	fs = sqlbuilder.Fields{
		p.UserIdField,
		p.PasswordField,
	}
	return fs
}

// IsEqual 判断明文密码是否和加密密码一致
func IsEqual(userPassword string, dbPassword string) bool {
	ok := strings.EqualFold(dbPassword, EncodingFn(userPassword))
	return ok
}

var EncodingFn = GetMD5Hash

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

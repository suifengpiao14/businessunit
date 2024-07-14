package password

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/businessunit/identity"
	"github.com/suifengpiao14/sqlbuilder"
)

func OptionPassword(f *sqlbuilder.Field) {
	f.SetName("password").SetTitle("密码").MergeSchema(sqlbuilder.Schema{
		Title:     "密码",
		Required:  true,
		Comment:   "对象标识,缺失时记录无意义",
		Type:      sqlbuilder.Schema_Type_string,
		MaxLength: 64,
		MinLength: 1,
	})
	f.WhereFns.InsertAsFirst(sqlbuilder.WhereValueFnShield) // 密码字段不会作为where条件
	f.ValueFns.AppendIfNotFirst(func(in any) (any, error) {
		password := cast.ToString(in)
		if password == "" {
			err := errors.Errorf("password request string")
			return nil, err
		}
		value := EncodingFn(password) // 对密码字段加密转换
		return value, nil
	})
}

func WithOptions(idField *sqlbuilder.Field, passwordField *sqlbuilder.Field) {
	idField.WithOptions(identity.OptionIdentity)
	passwordField.WithOptions(OptionPassword)
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

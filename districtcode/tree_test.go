package districtcode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/businessunit/districtcode"
)

type Record struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Children []Record `json:"children"`
}

func (r *Record) GetCode() string {
	return r.Code
}
func (r *Record) AddChildren(children ...districtcode.RecordI) {
	if len(r.Children) == 0 {
		r.Children = make([]Record, 0)
	}
	for _, child := range children {
		c := child.(*Record)
		r.Children = append(r.Children, *c)
	}

}

func Convert[T []districtcode.RecordI](records T) (recordIs []districtcode.RecordI) {
	recordIs = make([]districtcode.RecordI, 0)
	for _, r := range records {
		recordIs = append(recordIs, r)
	}
	return recordIs
}

type Records []*Record

func (rs Records) RecordIs() (recordIs []districtcode.RecordI) {
	recordIs = make([]districtcode.RecordI, 0)
	for _, r := range rs {
		recordIs = append(recordIs, r)
	}
	return recordIs
}

var records = Records{
	{
		Code: "440301",
		Name: "福田区",
	},
	{
		Code: "4403",
		Name: "深圳",
	},
	{
		Code: "44",
		Name: "广东",
	},
}

func TestTree(t *testing.T) {
	trees := districtcode.Tree(records.RecordIs())
	b, err := json.Marshal(trees)
	require.NoError(t, err)
	s := string(b)
	fmt.Println(s)
}

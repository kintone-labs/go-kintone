// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestSingleLineTextField(t *testing.T) {
	t.Parallel()

	var s SingleLineTextField = "hoge"
	if s != "hoge" {
		t.Fatal("hoge != hoge")
	}

	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "SINGLE_LINE_TEXT" {
		t.Error("Invalid type")
	}
	if m["value"] != "hoge" {
		t.Error("Invalid value")
	}
}

func TestMultiLineTextField(t *testing.T) {
	t.Parallel()

	var s MultiLineTextField = `hoge
fuga
`
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "MULTI_LINE_TEXT" {
		t.Error("Invalid type")
	}
	if m["value"] != `hoge
fuga
` {
		t.Error("Invalid value")
	}
}

func TestRichTextField(t *testing.T) {
	t.Parallel()

	var s RichTextField = `<a href="http://www.cybozu.com">hoge</a>
fuga
`
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "RICH_TEXT" {
		t.Error("Invalid type")
	}
	if m["value"] != `<a href="http://www.cybozu.com">hoge</a>
fuga
` {
		t.Error("Invalid value")
	}
}

func TestDecimalField(t *testing.T) {
	t.Parallel()

	var s DecimalField = "123456789012345678901234567890"
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "NUMBER" {
		t.Error("Invalid type")
	}
	if m["value"] != "123456789012345678901234567890" {
		t.Error("Invalid value")
	}
}

func TestCheckBoxField(t *testing.T) {
	t.Parallel()

	var s CheckBoxField = []string{"abc", "def"}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "CHECK_BOX" {
		t.Error("Invalid type")
	}
	if !reflect.DeepEqual(m["value"], []interface{}{"abc", "def"}) {
		t.Logf("%T", m["value"])
		t.Logf("%T", []string{"abc", "def"})
		t.Error("Invalid value")
	}
}

func TestRadioButtonField(t *testing.T) {
	t.Parallel()

	var s RadioButtonField = "button1"
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "RADIO_BUTTON" {
		t.Error("Invalid type")
	}
	if m["value"] != "button1" {
		t.Error("Invalid value")
	}
}

func TestSingleSelectField(t *testing.T) {
	t.Parallel()

	s := SingleSelectField{"select1", true}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "DROP_DOWN" {
		t.Error("Invalid type")
	}
	if m["value"] != "select1" {
		t.Error("Invalid value")
	}

	s2 := SingleSelectField{"", false}
	j, err = json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if m["value"] != nil {
		t.Error("value must be nil")
	}
}

func TestMultiSelectField(t *testing.T) {
	t.Parallel()

	var s MultiSelectField = []string{"select1", "select2"}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "MULTI_SELECT" {
		t.Error("Invalid type")
	}
	if !reflect.DeepEqual(m["value"], []interface{}{"select1", "select2"}) {
		t.Error("Invalid value")
	}
}

func TestFileField(t *testing.T) {
	t.Parallel()

	var s FileField = []File{
		{"text/plain", "12345", "aaa.txt", 12345678},
		{"application/octet-stream", "ghruu4", "hoge.wmv", 333},
	}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "FILE" {
		t.Error("Invalid type")
	}
	files := m["value"].([]interface{})
	if len(files) != 2 {
		t.Error("Wrong value length")
	}
	f1, ok := files[0].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid File data")
	}
	if _, ok := f1["contentType"]; !ok {
		t.Error("No contentType")
	}
	if f1["contentType"] != "text/plain" {
		t.Error("Invalid content type")
	}
	if _, ok := f1["fileKey"]; !ok {
		t.Error("No fileKey")
	}
	if f1["fileKey"] != "12345" {
		t.Error("Wrong file key")
	}
	if _, ok := f1["name"]; !ok {
		t.Error("No name")
	}
	if f1["name"] != "aaa.txt" {
		t.Error("Wrong file name")
	}
	if _, ok := f1["size"]; !ok {
		t.Error("No size")
	}
	if f1["size"] != "12345678" {
		t.Error("Wrong file size")
	}
}

func TestLinkField(t *testing.T) {
	t.Parallel()

	var s LinkField = "https://www.google.com/"
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "LINK" {
		t.Error("Invalid type")
	}
	if m["value"] != "https://www.google.com/" {
		t.Error("Invalid value")
	}
}

func TestDateField(t *testing.T) {
	t.Parallel()

	s := NewDateField(2000, time.January, 15)
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "DATE" {
		t.Error("Invalid type")
	}
	if m["value"] != "2000-01-15" {
		t.Error("Invalid value")
	}

	s2 := DateField{Valid: false}
	j, err = json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if m["value"] != nil {
		t.Error("value must be nil")
	}
}

func TestTimeField(t *testing.T) {
	t.Parallel()

	s := NewTimeField(19, 55)
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "TIME" {
		t.Error("Invalid type")
	}
	if m["value"] != "19:55:00" {
		t.Error("Invalid value")
	}
}

func TestDateTimeField(t *testing.T) {
	t.Parallel()

	s := NewDateTimeField(2014, time.February, 3, 9, 17)
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "DATETIME" {
		t.Error("Invalid type")
	}
	if m["value"] != "2014-02-03T09:17:00Z" {
		t.Error("Invalid value")
	}
}

func TestUserField(t *testing.T) {
	t.Parallel()

	var s UserField = []User{
		{"ymmt2005", "Yamamoto, Hirotaka"},
	}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "USER_SELECT" {
		t.Error("Invalid type")
	}
	users := m["value"].([]interface{})
	if len(users) != 1 {
		t.Error("Wrong value length")
	}
	u1, ok := users[0].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid User data")
	}
	if _, ok := u1["code"]; !ok {
		t.Error("No code")
	}
	if u1["code"] != "ymmt2005" {
		t.Error("Invalid code")
	}
	if _, ok := u1["name"]; !ok {
		t.Error("No name")
	}
	if u1["name"] != "Yamamoto, Hirotaka" {
		t.Error("Wrong name")
	}
}

func TestCategoryField(t *testing.T) {
	t.Parallel()

	var s CategoryField = []string{"cat1", "cat2", "cat3"}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "CATEGORY" {
		t.Error("Invalid type")
	}
	if !reflect.DeepEqual(m["value"], []interface{}{"cat1", "cat2", "cat3"}) {
		t.Error("Invalid value")
	}
}

func TestStatusField(t *testing.T) {
	t.Parallel()

	var s StatusField = "status"
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "STATUS" {
		t.Error("Invalid type")
	}
	if m["value"] != "status" {
		t.Error("Invalid value")
	}
}

func TestAssigneeField(t *testing.T) {
	t.Parallel()

	var s AssigneeField = []User{
		{"ymmt2005", "Yamamoto, Hirotaka"},
		{"foobar", "zot"},
	}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "STATUS_ASSIGNEE" {
		t.Error("Invalid type")
	}
	users := m["value"].([]interface{})
	if len(users) != 2 {
		t.Error("Wrong value length")
	}
	u1, ok := users[1].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid User data")
	}
	if _, ok := u1["code"]; !ok {
		t.Error("No code")
	}
	if u1["code"] != "foobar" {
		t.Error("Invalid code")
	}
	if _, ok := u1["name"]; !ok {
		t.Error("No name")
	}
	if u1["name"] != "zot" {
		t.Error("Wrong name")
	}
}

func TestRecordNumberField(t *testing.T) {
	t.Parallel()

	var s RecordNumberField = "12345"
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "RECORD_NUMBER" {
		t.Error("Invalid type")
	}
	if m["value"] != "12345" {
		t.Error("Invalid value")
	}
}

func TestCreatorField(t *testing.T) {
	t.Parallel()

	s := CreatorField{"ymmt2005", "Yamamoto, Hirotaka"}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "CREATOR" {
		t.Error("Invalid type")
	}
	user, ok := m["value"].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid User data")
	}
	if _, ok := user["code"]; !ok {
		t.Error("No code")
	}
	if user["code"] != "ymmt2005" {
		t.Error("Invalid code")
	}
	if _, ok := user["name"]; !ok {
		t.Error("No name")
	}
	if user["name"] != "Yamamoto, Hirotaka" {
		t.Error("Wrong name")
	}
}

func TestModifierField(t *testing.T) {
	t.Parallel()

	s := ModifierField{"ymmt2005", "Yamamoto, Hirotaka"}
	j, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["type"]; !ok {
		t.Error("Invalid marshal data!")
	}
	if m["type"] != "MODIFIER" {
		t.Error("Invalid type")
	}
	user, ok := m["value"].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid User data")
	}
	if _, ok := user["code"]; !ok {
		t.Error("No code")
	}
	if user["code"] != "ymmt2005" {
		t.Error("Invalid code")
	}
	if _, ok := user["name"]; !ok {
		t.Error("No name")
	}
	if user["name"] != "Yamamoto, Hirotaka" {
		t.Error("Wrong name")
	}
}

func TestSubTableField(t *testing.T) {
	t.Parallel()

	s := []SubTableEntry{
		{Id: "123", Value: map[string]interface{}{
			"abc": RecordNumberField("12345")}}}
	t.Log(s)
}

func TestIsBuiltinField(t *testing.T) {
	t.Parallel()

	if IsBuiltinField(SingleLineTextField("aaa")) {
		t.Error("SingleLineTextField is not built-in")
	}
	if !IsBuiltinField(CreatorField{"user1", "user name"}) {
		t.Error("CreatorField is built-in")
	}
}

// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func newApp(appID uint64) *App {
	return &App{
		Domain:   os.Getenv("KINTONE_DOMAIN"),
		User:     os.Getenv("KINTONE_USER"),
		Password: os.Getenv("KINTONE_PASSWORD"),
		AppId:    appID,
	}
}

func newAppWithApiToken(appId uint64) *App {
	return &App{
		Domain:   os.Getenv("KINTONE_DOMAIN"),
		ApiToken: os.Getenv("KINTONE_API_TOKEN"),
		AppId:    appId,
	}
}

func newAppInGuestSpace(appId uint64, guestSpaceId uint64) *App {
	return &App{
		Domain:       os.Getenv("KINTONE_DOMAIN"),
		User:         os.Getenv("KINTONE_USER"),
		Password:     os.Getenv("KINTONE_PASSWORD"),
		AppId:        appId,
		GuestSpaceId: guestSpaceId,
	}
}

func TestGetRecord(t *testing.T) {
	a := newApp(4799)
	if len(a.Password) == 0 {
		t.Skip()
	}

	if rec, err := a.GetRecord(116); err != nil {
		t.Error(err)
	} else {
		if rec.Id() != 116 {
			t.Errorf("Unexpected Id: %d", rec.Id())
		}
		for _, f := range rec.Fields {
			if files, ok := f.(FileField); ok {
				if len(files) == 0 {
					continue
				}
				fd, err := a.Download(files[0].FileKey)
				if err != nil {
					t.Error(err)
				} else {
					data, _ := ioutil.ReadAll(fd.Reader)
					t.Logf("%s %d bytes", fd.ContentType, len(data))
				}
			}
		}
	}

	if recs, err := a.GetRecords(nil, "limit 3 offset 3"); err != nil {
		t.Error(err)
	} else {
		if len(recs) > 3 {
			t.Error("Too many records")
		}
	}

	if recs, err := a.GetAllRecords([]string{"レコード番号"}); err != nil {
		t.Error(err)
	} else {
		t.Log(len(recs))
	}
}
func TestGetAllRecord(t *testing.T) {
	os.Setenv("KINTONE_DOMAIN", "tgbn0.kintone.com")
	os.Setenv("KINTONE_USER", "Administrator")
	os.Setenv("KINTONE_PASSWORD", "cybozu123")
	a := newApp(18)
	result, _ := a.GetAllRecords([]string{"id"})
	fmt.Println(result)
}
func TestAddRecord(t *testing.T) {
	a := newApp(9004)
	if len(a.Password) == 0 {
		t.Skip()
	}

	fileKey, err := a.Upload("ほげ春巻.txta", "text/html",
		bytes.NewReader([]byte(`abc
<a href="https://www.cybozu.com/">hoge</a>
`)))
	if err != nil {
		t.Error("Upload failed", err)
	}

	rec := NewRecord(map[string]interface{}{
		"title": SingleLineTextField("test!"),
		"file": FileField{
			{FileKey: fileKey},
		},
	})
	_, err = a.AddRecord(rec)
	if err != nil {
		t.Error("AddRecord failed", rec)
	}

	recs := []*Record{
		NewRecord(map[string]interface{}{
			"title": SingleLineTextField("multi add 1"),
		}),
		NewRecord(map[string]interface{}{
			"title": SingleLineTextField("multi add 2"),
		}),
	}
	ids, err := a.AddRecords(recs)
	if err != nil {
		t.Error("AddRecords failed", recs)
	} else {
		t.Log(ids)
	}
}

func TestUpdateRecord(t *testing.T) {
	a := newApp(9004)
	if len(a.Password) == 0 {
		t.Skip()
	}

	rec, err := a.GetRecord(4)
	if err != nil {
		t.Fatal(err)
	}
	rec.Fields["title"] = SingleLineTextField("new title")
	if err := a.UpdateRecord(rec, true); err != nil {
		t.Error("UpdateRecord failed", err)
	}

	if err := a.UpdateRecordByKey(rec, true, "key"); err != nil {
		t.Error("UpdateRecordByKey failed", err)
	}

	recs, err := a.GetRecords(nil, "limit 3")
	if err != nil {
		t.Fatal(err)
	}
	for _, rec := range recs {
		rec.Fields["title"] = SingleLineTextField(time.Now().String())
	}
	if err := a.UpdateRecords(recs, true); err != nil {
		t.Error("UpdateRecords failed", err)
	}

	if err := a.UpdateRecordsByKey(recs, true, "key"); err != nil {
		t.Error("UpdateRecordsByKey failed", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	a := newApp(9004)
	if len(a.Password) == 0 {
		t.Skip()
	}

	ids := []uint64{6, 7}
	if err := a.DeleteRecords(ids); err != nil {
		t.Error("DeleteRecords failed", err)
	}
}

func TestGetCursor(t *testing.T) {
	app := newApp(18)
	result, err := app.getCurSor("2e17d05f-d046-494e-ae8a-e2f82d68ff05")
	if len(app.Password) == 0 {
		t.Skip()
	}
	if err != nil {
		t.Errorf("get cursor is fail: %v", err)
	}
	fmt.Println(string(result))

}
func TestDeleteCursor(t *testing.T) {
	app := newApp(18)
	if len(app.Password) == 0 {
		t.Skip()
	}
	result, err := app.deleteCursor("b997c4d6-b514-4567-9c00-1014770ead5c")
	if err != nil {
		t.Errorf("delete cursor is fail: %v", err)
	}
	fmt.Println(result)
}
func TestCreateCurSor(t *testing.T) {
	app := newAppWithApiToken(18)
	if len(app.ApiToken) == 0 {
		t.Skip()
	}
	result, err := app.createCursor([]string{"$id", "date"})
	if err != nil {
		t.Errorf("create cursor is fail: %v", err)
	}
	fmt.Println("result", string(result))
}

func TestFields(t *testing.T) {
	a := newApp(8326)
	if len(a.Password) == 0 {
		t.Skip()
	}

	fi, err := a.Fields()
	if err != nil {
		t.Error("Fields failed", err)
	}
	for _, f := range fi {
		t.Log(f)
	}
}

func TestApiToken(t *testing.T) {
	a := newAppWithApiToken(9974)
	if len(a.ApiToken) == 0 {
		t.Skip()
	}

	_, err := a.Fields()
	if err != nil {
		t.Error("Api token failed", err)
	}
}

func TestGuestSpace(t *testing.T) {
	a := newAppInGuestSpace(185, 9)
	if len(a.Password) == 0 {
		t.Skip()
	}

	_, err := a.Fields()
	if err != nil {
		t.Error("GuestSpace failed", err)
	}
}

func TestGetRecordComments(t *testing.T) {
	a := newApp(13)
	var offset uint64 = 5
	var limit uint64 = 10
	if rec, err := a.GetRecordComments(3, "asc", offset, limit); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(rec[0].Id, "6") {
			t.Errorf("the first comment id mismatch. expected is 6 but actual %v", rec[0].Id)
		}
	}
}
func TestAddRecordComment(t *testing.T) {
	appTest := newApp(12)
	mentionMemberCybozu := &ObjMention{Code: "cybozu", Type: ConstCommentMentionTypeUser}
	mentionGroupAdmin := &ObjMention{Code: "Administrators", Type: ConstCommentMentionTypeGroup}
	mentionDepartmentAdmin := &ObjMention{Code: "Admin", Type: ConstCommentMentionTypeDepartment}
	var cmt Comment
	cmt.Text = "Test comment 222"
	cmt.Mentions = []*ObjMention{mentionGroupAdmin, mentionMemberCybozu, mentionDepartmentAdmin}
	cmtID, err := appTest.AddRecordComment(2, &cmt)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("return value(comment-id) is %v", cmtID)
	}
}

func TestDeleteComment(t *testing.T) {
	appTest := newApp(4)
	var cmtID uint64 = 14
	err := appTest.DeleteComment(3, 12)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("The comment with id =  %v has been deleted successefully!", cmtID)
	}
}

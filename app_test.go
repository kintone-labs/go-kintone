// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
	"strconv"
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

func TestAddRecord(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
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
		"key":DecimalField("1"),
		"file": FileField{
			{FileKey: fileKey},
		},
	})
	id, err := a.AddRecord(rec)
	if err != nil {
		t.Error("AddRecord failed", err)
	}
	os.Setenv("KINTONE_TEST_REC_ID", id)

	recs := []*Record{
		NewRecord(map[string]interface{}{
			"key":DecimalField("2"),
			"title": SingleLineTextField("multi add 1"),
		}),
		NewRecord(map[string]interface{}{
			"key":DecimalField("3"),
			"title": SingleLineTextField("multi add 2"),
		}),
	}
	ids, err := a.AddRecords(recs)
	if err != nil {
		t.Error("AddRecords failed", err)
	} else {
		t.Log(ids)
	}
}

func TestGetRecord(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
	if len(a.Password) == 0 {
		t.Skip()
	}

	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	if rec, err := a.GetRecord(recId); err != nil {
		t.Error(err)
	} else {
		if rec.Id() != recId {
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

func TestUpdateRecord(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
	if len(a.Password) == 0 {
		t.Skip()
	}

	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	rec, err := a.GetRecord(recId)
	if err != nil {
		t.Fatal(err)
	}
	rec.Fields["title"] = SingleLineTextField("new title")
	delete(rec.Fields, "レコード番号")
	delete(rec.Fields, "作成者")
	delete(rec.Fields, "更新者")
	delete(rec.Fields, "作成日時")
	delete(rec.Fields, "更新日時")
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
		delete(rec.Fields, "レコード番号")
		delete(rec.Fields, "作成者")
		delete(rec.Fields, "更新者")
		delete(rec.Fields, "作成日時")
		delete(rec.Fields, "更新日時")
	}
	if err := a.UpdateRecords(recs, true); err != nil {
		t.Error("UpdateRecords failed", err)
	}

	if err := a.UpdateRecordsByKey(recs, true, "key"); err != nil {
		t.Error("UpdateRecordsByKey failed", err)
	}
}

func TestFields(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
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
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newAppWithApiToken(appId)
	if len(a.ApiToken) == 0 {
		t.Skip()
	}

	_, err := a.Fields()
	if err != nil {
		t.Error("Api token failed", err)
	}
}

func TestGuestSpace(t *testing.T) {
	guestAppId, _ := strconv.ParseUint(os.Getenv("KINTONE_GUEST_APP_ID"), 10, 64)
	guestSpaceId, _ := strconv.ParseUint(os.Getenv("KINTONE_GUEST_SPACE_ID"), 10, 64)
	a := newAppInGuestSpace(guestAppId, guestSpaceId)
	if len(a.Password) == 0 {
		t.Skip()
	}

	_, err := a.Fields()
	if err != nil {
		t.Error("GuestSpace failed", err)
	}
}

func TestAddRecordComment(t *testing.T) {
	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	appTest := newApp(appId)
	mentionMemberUser := &ObjMention{Code: os.Getenv("KINTONE_USER"), Type: ConstCommentMentionTypeUser}
	mentionGroupAdmin := &ObjMention{Code: "Administrators", Type: ConstCommentMentionTypeGroup}
	mentionDepartmentAdmin := &ObjMention{Code: "Admin", Type: ConstCommentMentionTypeDepartment}
	var cmt Comment
	cmt.Text = "Test comment 222"
	cmt.Mentions = []*ObjMention{mentionGroupAdmin, mentionMemberUser, mentionDepartmentAdmin}

	for i := 0; i < 20; i++ {
		cmtID, err := appTest.AddRecordComment(recId, &cmt)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("return value(comment-id) is %v", cmtID)
		}
	}
}

func TestGetRecordComments(t *testing.T) {
	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
	var offset uint64 = 5
	var limit uint64 = 10
	if rec, err := a.GetRecordComments(recId, "asc", offset, limit); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(rec[0].Id, "6") {
			t.Errorf("the first comment id mismatch. expected is 6 but actual %v", rec[0].Id)
		}
	}
}

func TestDeleteComment(t *testing.T) {
	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	appTest := newApp(appId)
	var cmtID uint64 = 14
	err := appTest.DeleteComment(recId, cmtID)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("The comment with id =  %v has been deleted successfully!", cmtID)
	}
}

func TestOpenCursor(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	appTest := newApp(appId)
	cursor, err := appTest.OpenCursor(nil, "key >= 2", 1)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("The cursor with id:%v, totalCount:%v is opened successfully!", cursor.Id, cursor.TotalCount)
		os.Setenv("KINTONE_CURSOR_ID", cursor.Id)
	}
}

func TestReadCursor(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	appTest := newApp(appId)
	if recs, err := appTest.ReadCursor(os.Getenv("KINTONE_CURSOR_ID")); err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", recs)
		if len(recs) > 1 {
			t.Error("Too many records")
		}
	}
}

func TestCloseCursor(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	appTest := newApp(appId)
	err := appTest.CloseCursor(os.Getenv("KINTONE_CURSOR_ID"))

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Cursor id:%v is closed successfully!", os.Getenv("KINTONE_CURSOR_ID"))
	}
}

func TestDeleteRecord(t *testing.T) {
	appId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_APP_ID"), 10, 64)
	a := newApp(appId)
	if len(a.Password) == 0 {
		t.Skip()
	}

	recId, _ := strconv.ParseUint(os.Getenv("KINTONE_TEST_REC_ID"), 10, 64)
	ids := []uint64{recId, recId+1, recId+2}
	if err := a.DeleteRecords(ids); err != nil {
		t.Error("DeleteRecords failed", err)
	}
}
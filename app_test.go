// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func newAppForTest(domain string, user string, password string, appID uint64) *App {
	return &App{
		Domain:   domain,
		User:     user,
		Password: password,
		AppId:    appID,
	}
}

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

func createResponseLocalTestServer(data string) (*httptest.Server, error) {
	ts, err := NewLocalHTTPSTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, data)
	}))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func NewLocalHTTPSTestServer(handler http.Handler) (*httptest.Server, error) {
	ts := httptest.NewUnstartedServer(handler)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	l, err := net.Listen("tcp", "localhost:8088")
	if err != nil {
		fmt.Println(err)
	}

	ts.Listener.Close()
	ts.Listener = l
	ts.StartTLS()

	return ts, nil
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

	json := `{
    "record": {
			"Updated_by": {
				"type": "MODIFIER",
				"value": {
						"code": "Administrator",
						"name": "Administrator"
				}
		},
        "$id": {
            "type": "__ID__",
            "value": "1"
        }
    }
}
`
	ts, _ := createResponseLocalTestServer(json)
	defer ts.Close()
	a := newAppForTest("127.0.0.1:8088", "test", "test", 2)

	if rec, err := a.GetRecord(1); err != nil {
		t.Error(err)
	} else {
		if rec.Id() != 1 {
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
	json := `
	{
		"comments": [
        {
            "id": "3",
            "text": "user14 Thank you! Looks great.",
            "createdAt": "2016-05-09T18:29:05Z",
            "creator": {
                "code": "user13",
                "name": "user13"
            },
            "mentions": [
                {
                    "code": "user14",
                    "type": "USER"
                }
            ]
        },
        {
            "id": "2",
            "text": "user13 Global Sales APAC Taskforce \nHere is today's report.",
            "createdAt": "2016-05-09T18:27:54Z",
            "creator": {
                "code": "user14",
                "name": "user14"
            },
            "mentions": [
                {
                    "code": "user13",
                    "type": "USER"
                },
                {
                    "code": "Global Sales_1BNZeQ",
                    "type": "ORGANIZATION"
                },
                {
                    "code": "APAC Taskforce_DJrvzu",
                    "type": "GROUP"
                }
            ]
        }
    ],
    "older": false,
    "newer": false
	}
	`
	ts, _ := createResponseLocalTestServer(json)
	defer ts.Close()
	a := newAppForTest("127.0.0.1:8088", "test", "test", 2)
	var offset uint64 = 0
	var limit uint64 = 10
	if rec, err := a.GetRecordComments(1, "asc", offset, limit); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(rec[0].Id, "3") {
			t.Errorf("the first comment id mismatch. expected is 3 but actual %v", rec[0].Id)
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

// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	KINTONE_DOMAIN         = "localhost:8088"
	KINTONE_USERNAME       = "test"
	KINTONE_PASSWORD       = "test"
	KINTONE_APP_ID         = 1
	KINTONE_API_TOKEN      = "1e42da75-8432-4adb-9a2b-dbb6e7cb3c6b"
	KINTONE_GUEST_SPACE_ID = 1
)

func createServerTest(mux *http.ServeMux) (*httptest.Server, error) {
	ts := httptest.NewUnstartedServer(mux)
	listen, err := net.Listen("tcp", KINTONE_DOMAIN)

	if err != nil {
		return nil, err
	}

	ts.Listener.Close()
	ts.Listener = listen
	ts.StartTLS()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return ts, nil
}

func createServerMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/k/v1/record.json", handleResponseGetRecord)
	mux.HandleFunc("/k/v1/records.json", handleResponseGetRecords)
	mux.HandleFunc("/k/v1/record/comments.json", handleResponseGetRecordsComments)
	mux.HandleFunc("/k/v1/file.json", handleResponseUploadFile)
	mux.HandleFunc("/k/v1/record/comment.json", handleResponseRecordComments)
	mux.HandleFunc("/k/v1/records/cursor.json", handleResponseRecordsCursor)
	mux.HandleFunc("/k/v1/form.json", handleResponseForm)
	mux.HandleFunc("/k/guest/1/v1/form.json", handleResponseForm)
	return mux, nil
}

// handler mux
func handleResponseForm(response http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		testData := GetDataTestForm()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseRecordsCursor(response http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		testData := GetDataTestGetRecordsByCursor()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		testData := GetTestDataDeleteCursor()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "POST" {
		testData := GetTestDataCreateCursor()
		fmt.Fprint(response, testData.output)
	}
}
func handleResponseRecordComments(response http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		testData := GetTestDataAddRecordComment()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		testData := GetDataTestDeleteRecordComment()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseUploadFile(response http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		testData := GetDataTestUploadFile()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseGetRecord(response http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		testData := GetTestDataGetRecord()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "PUT" {
		testData := GetTestDataUpdateRecordByKey()
		fmt.Fprint(response, testData.output)
	}

}

func handleResponseGetRecords(response http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		testData := GetTestDataGetRecords()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		testData := GetTestDataDeleteRecords()
		fmt.Fprint(response, testData.output)
	}

}

func handleResponseGetRecordsComments(response http.ResponseWriter, r *http.Request) {
	testData := GetDataTestRecordComments()
	fmt.Fprint(response, testData.output)

}

func TestMain(m *testing.M) {
	mux, err := createServerMux()
	if err != nil {
		fmt.Println("StartServerTest", err)
	}
	ts, err := createServerTest(mux)
	if err != nil {
		fmt.Println("createServerTest", err)
	}
	m.Run()
	ts.Close()
}

func newApp() *App {
	return &App{
		Domain:   KINTONE_DOMAIN,
		User:     KINTONE_USERNAME,
		Password: KINTONE_PASSWORD,
		AppId:    KINTONE_APP_ID,
		ApiToken: KINTONE_API_TOKEN,
	}
}
func newAppWithGuest() *App {
	return &App{
		Domain:       KINTONE_DOMAIN,
		AppId:        KINTONE_APP_ID,
		ApiToken:     KINTONE_API_TOKEN,
		GuestSpaceId: KINTONE_GUEST_SPACE_ID,
	}
}
func newAppWithToken() *App {
	return &App{
		Domain:   KINTONE_DOMAIN,
		ApiToken: KINTONE_API_TOKEN,
	}
}

func TestAddRecord(t *testing.T) {
	testData := GetDataTestAddRecord()
	a := newApp()

	fileKey, err := a.Upload(testData.input[0].(string), "text/html",
		testData.input[1].(io.Reader))
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
func TestGetRecord(t *testing.T) {
	testData := GetTestDataGetRecord()
	a := newApp()
	if rec, err := a.GetRecord(uint64(testData.input[0].(int))); err != nil {
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
func TestUpdateRecord(t *testing.T) {
	testData := GetTestDataGetRecord()
	a := newApp()

	rec, err := a.GetRecord(uint64(testData.input[0].(int)))
	if err != nil {
		t.Fatal(err)
	}
	rec.Fields["title"] = SingleLineTextField("new title")
	if err := a.UpdateRecord(rec, true); err != nil {
		t.Error("UpdateRecord failed", err)
	}

	rec.Fields["key"] = SingleLineTextField(` {
		"field": "unique_key",
		"value": "unique_code"
	}`)
	if err := a.UpdateRecordByKey(rec, true, "key"); err != nil {

		t.Error("UpdateRecordByKey failed", err)
	}
	recs, err := a.GetRecords(nil, "limit 3")
	if err != nil {
		t.Fatal(err)
	}

	for _, rec := range recs {
		rec.Fields["title"] = SingleLineTextField(time.Now().String())
		rec.Fields["key"] = SingleLineTextField(` {
			"field": "unique_key",
			"value": "unique_code"
	}`)
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

func TestGetRecordsByCursor(t *testing.T) {
	app := newApp(18)

	if len(app.Password) == 0 {
		t.Skip()
	}

	cursor := app.createCursorForTest()
	record, err := app.GetRecordsByCursor(string(cursor.Id))

	if err != nil {
		t.Errorf("TestGetCursor is failed: %v", err)
	}
	fmt.Println(record)

}

func (app *App) createCursorForTest() *Cursor {
	cursor, err := app.CreateCursor([]string{"$id", "Status"}, "", 400)
	fmt.Println("cursor", cursor)
	if err != nil {
		fmt.Println("createCursorForTest failed: ", err)
	}
	return cursor
}

func TestDeleteCursor(t *testing.T) {
	app := newApp(18)
	if len(app.Password) == 0 {
		t.Skip()
	}

	cursor := app.createCursorForTest()
	fmt.Println("cursor", cursor)
	err := app.DeleteCursor(string(cursor.Id))

	if err != nil {
		t.Errorf("TestDeleteCursor is failed: %v", err)
	}
}

func TestCreateCursor(t *testing.T) {
	app := newApp(18)
	if len(app.Password) == 0 {
		t.Skip()
	}
	_, err := app.CreateCursor([]string{"$id", "date"}, "", 100)
	if err != nil {
		t.Errorf("TestCreateCurSor is failed: %v", err)
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
	a := newAppForTest()
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

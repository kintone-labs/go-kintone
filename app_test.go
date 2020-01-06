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
	AUTH_HEADER_TOKEN      = "X-Cybozu-API-Token"
	AUTH_HEADER_PASSWORD   = "X-Cybozu-Authorization"
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

// header check
func checkAuth(response http.ResponseWriter, r *http.Request) {
	authPassword := r.Header.Get(AUTH_HEADER_PASSWORD)
	authToken := r.Header.Get(AUTH_HEADER_TOKEN)
	if authPassword == "" && authToken == "" {
		http.Error(response, "Authenticated is fail!", http.StatusForbidden)

		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("500 - Authenticated is fail!"))
	}
}
func checkContentType(response http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		fmt.Println("Content-type is invalid")
		http.Error(response, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	}
}

// handler mux
func handleResponseForm(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "GET" {
		checkContentType(response, r)
		testData := GetDataTestForm()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseRecordsCursor(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "GET" {
		testData := GetDataTestGetRecordsByCursor()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		checkContentType(response, r)
		testData := GetTestDataDeleteCursor()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "POST" {
		checkContentType(response, r)
		testData := GetTestDataCreateCursor()
		fmt.Fprint(response, testData.output)
	}
}
func handleResponseRecordComments(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "POST" {
		checkContentType(response, r)
		testData := GetTestDataAddRecordComment()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		checkContentType(response, r)
		testData := GetDataTestDeleteRecordComment()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseUploadFile(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "POST" {
		testData := GetDataTestUploadFile()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseGetRecord(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "GET" {
		checkContentType(response, r)
		testData := GetTestDataGetRecord()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "PUT" {
		checkContentType(response, r)
		testData := GetTestDataUpdateRecordByKey()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "POST" {
		checkContentType(response, r)
		testData := GetTestDataAddRecord()
		fmt.Fprint(response, testData.output)
	}

}

func handleResponseGetRecords(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	if r.Method == "GET" {
		checkContentType(response, r)
		testData := GetTestDataGetRecords()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "DELETE" {
		checkContentType(response, r)
		testData := GetTestDataDeleteRecords()
		fmt.Fprint(response, testData.output)
	} else if r.Method == "POST" {
		checkContentType(response, r)
		testData := GetTestDataAddRecords()
		fmt.Fprint(response, testData.output)
	}

}

func handleResponseGetRecordsComments(response http.ResponseWriter, r *http.Request) {
	checkAuth(response, r)
	checkContentType(response, r)
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
	a := newApp()

	ids := []uint64{6, 7}
	if err := a.DeleteRecords(ids); err != nil {
		t.Error("DeleteRecords failed", err)
	}
}

func TestGetRecordsByCursor(t *testing.T) {
	testData := GetDataTestGetRecordsByCursor()
	app := newApp()
	_, err := app.GetRecordsByCursor(testData.input[0].(string))
	if err != nil {
		t.Errorf("TestGetCursor is failed: %v", err)
	}

}

func TestDeleteCursor(t *testing.T) {
	testData := GetTestDataDeleteCursor()
	app := newApp()
	err := app.DeleteCursor(testData.input[0].(string))
	if err != nil {
		t.Errorf("TestDeleteCursor is failed: %v", err)
	}
}

func TestCreateCursor(t *testing.T) {
	testData := GetTestDataCreateCursor()
	app := newApp()
	_, err := app.CreateCursor(testData.input[0].([]string), testData.input[1].(string), uint64(testData.input[2].(int)))
	if err != nil {
		t.Errorf("TestCreateCurSor is failed: %v", err)
	}
}

func TestFields(t *testing.T) {
	a := newApp()

	fi, err := a.Fields()
	if err != nil {
		t.Error("Fields failed", err)
	}
	for _, f := range fi {
		t.Log(f)
	}
}

func TestApiToken(t *testing.T) {
	a := newAppWithToken()
	_, err := a.Fields()
	if err != nil {
		t.Error("Api token failed", err)
	}
}

func TestGuestSpace(t *testing.T) {
	a := newAppWithGuest()

	_, err := a.Fields()
	if err != nil {
		t.Error("GuestSpace failed", err)
	}
}

func TestGetRecordComments(t *testing.T) {
	a := newApp()
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
	testData := GetTestDataAddRecordComment()
	appTest := newApp()
	mentionMemberCybozu := &ObjMention{Code: "cybozu", Type: ConstCommentMentionTypeUser}
	mentionGroupAdmin := &ObjMention{Code: "Administrators", Type: ConstCommentMentionTypeGroup}
	mentionDepartmentAdmin := &ObjMention{Code: "Admin", Type: ConstCommentMentionTypeDepartment}
	var cmt Comment
	cmt.Text = "Test comment 222"
	cmt.Mentions = []*ObjMention{mentionGroupAdmin, mentionMemberCybozu, mentionDepartmentAdmin}
	cmtID, err := appTest.AddRecordComment(uint64(testData.input[0].(int)), &cmt)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("return value(comment-id) is %v", cmtID)
	}
}

func TestDeleteComment(t *testing.T) {
	appTest := newApp()
	var cmtID uint64 = 14
	err := appTest.DeleteComment(3, 12)

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("The comment with id =  %v has been deleted successefully!", cmtID)
	}
}

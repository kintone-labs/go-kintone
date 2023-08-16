// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
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
	AUTH_HEADER_BASIC      = "Authorization"
	CONTENT_TYPE           = "Content-Type"
	APPLICATION_JSON       = "application/json"
	BASIC_AUTH             = true
	BASIC_AUTH_USER        = "basic"
	BASIC_AUTH_PASSWORD    = "basic"
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

func createServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/k/v1/record.json", handleResponseGetRecord)
	mux.HandleFunc("/k/v1/records.json", handleResponseGetRecords)
	mux.HandleFunc("/k/v1/record/comments.json", handleResponseGetRecordsComments)
	mux.HandleFunc("/k/v1/file.json", handleResponseUploadFile)
	mux.HandleFunc("/k/v1/record/comment.json", handleResponseRecordComments)
	mux.HandleFunc("/k/v1/records/cursor.json", handleResponseRecordsCursor)
	mux.HandleFunc("/k/v1/app/status.json", handleResponseProcess)
	mux.HandleFunc("/k/v1/form.json", handleResponseForm)
	mux.HandleFunc("/k/guest/1/v1/form.json", handleResponseForm)
	return mux
}

// header check
func checkAuth(response http.ResponseWriter, request *http.Request) {
	authPassword := request.Header.Get(AUTH_HEADER_PASSWORD)
	authToken := request.Header.Get(AUTH_HEADER_TOKEN)
	authBasic := request.Header.Get(AUTH_HEADER_BASIC)

	userAndPass := base64.StdEncoding.EncodeToString(
		[]byte(KINTONE_USERNAME + ":" + KINTONE_PASSWORD))

	userAndPassBasic := "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(BASIC_AUTH_USER+":"+BASIC_AUTH_PASSWORD))

	if authToken == "" && authPassword == "" && authBasic == "" {
		http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	if BASIC_AUTH && authBasic != "" && authBasic != userAndPassBasic {
		http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	if authToken != "" && authToken != KINTONE_API_TOKEN {
		http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	if authPassword != "" && authPassword != userAndPass {
		http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func checkContentType(response http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get(CONTENT_TYPE)
	if contentType != APPLICATION_JSON {
		http.Error(response, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}
}

// handler mux
func handleResponseProcess(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	TestData := GetTestDataProcess()
	fmt.Fprint(response, TestData.output)
}

func handleResponseForm(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	if request.Method == "GET" {
		checkContentType(response, request)
		testData := GetDataTestForm()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseRecordsCursor(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	if request.Method == "GET" {
		testData := GetDataTestGetRecordsByCursor()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "DELETE" {
		checkContentType(response, request)
		testData := GetTestDataDeleteCursor()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "POST" {
		checkContentType(response, request)
		testData := GetTestDataCreateCursor()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseRecordComments(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	checkContentType(response, request)
	if request.Method == "POST" {
		testData := GetTestDataAddRecordComment()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "DELETE" {
		testData := GetDataTestDeleteRecordComment()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseUploadFile(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	if request.Method == "POST" {
		testData := GetDataTestUploadFile()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseGetRecord(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	checkContentType(response, request)
	if request.Method == "GET" {
		testData := GetTestDataGetRecord()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "PUT" {
		testData := GetTestDataUpdateRecordByKey()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "POST" {
		testData := GetTestDataAddRecord()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseGetRecords(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	checkContentType(response, request)
	if request.Method == "GET" {
		type RequestBody struct {
			App        uint64   `json:"app,string"`
			Fields     []string `json:"fields"`
			Query      string   `json:"query"`
			TotalCount bool     `json:"totalCount"`
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(response, "Bad request", http.StatusBadRequest)
			return
		}
		var bodyRequest RequestBody
		if err := json.Unmarshal([]byte(body), &bodyRequest); err != nil {
			http.Error(response, "Body incorrect", http.StatusBadRequest)
		}

		if bodyRequest.TotalCount {
			testData := GetTestDataGetRecordsWithTotalCount()
			fmt.Fprint(response, testData.output)
		} else {
			testData := GetTestDataGetRecords()
			fmt.Fprint(response, testData.output)
		}
	} else if request.Method == "DELETE" {
		testData := GetTestDataDeleteRecords()
		fmt.Fprint(response, testData.output)
	} else if request.Method == "POST" {
		testData := GetTestDataAddRecords()
		fmt.Fprint(response, testData.output)
	}
}

func handleResponseGetRecordsComments(response http.ResponseWriter, request *http.Request) {
	checkAuth(response, request)
	checkContentType(response, request)
	testData := GetDataTestRecordComments()
	fmt.Fprint(response, testData.output)
}

func TestMain(m *testing.M) {
	mux := createServerMux()
	ts, err := createServerTest(mux)
	if err != nil {
		fmt.Println("createServerTest: ", err)
		os.Exit(1)
	}
	code := m.Run()
	ts.Close()
	os.Exit(code)
}

func newApp() *App {
	return &App{
		Domain:   KINTONE_DOMAIN,
		User:     KINTONE_USERNAME,
		Password: KINTONE_PASSWORD,
		AppId:    KINTONE_APP_ID,
	}
}

func newAppWithGuest() *App {
	return &App{
		Domain:       KINTONE_DOMAIN,
		AppId:        KINTONE_APP_ID,
		User:         KINTONE_USERNAME,
		Password:     KINTONE_PASSWORD,
		GuestSpaceId: KINTONE_GUEST_SPACE_ID,
	}
}

func newAppWithToken() *App {
	return &App{
		AppId:    KINTONE_APP_ID,
		Domain:   KINTONE_DOMAIN,
		ApiToken: KINTONE_API_TOKEN,
	}
}

func TestAddRecord(t *testing.T) {
	testData := GetDataTestAddRecord()
	app := newApp()

	fileKey, err := app.Upload(testData.input[0].(string), testData.input[2].(string),
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
	_, err = app.AddRecord(rec)
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
	ids, err := app.AddRecords(recs)
	if err != nil {
		t.Error("AddRecords failed", recs)
	} else {
		t.Log(ids)
	}
}

func TestGetRecord(t *testing.T) {
	testData := GetTestDataGetRecord()
	testDataRecords := GetTestDataGetRecords()
	app := newApp()
	if rec, err := app.GetRecord(uint64(testData.input[0].(int))); err != nil {
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
				fd, err := app.Download(files[0].FileKey)
				if err != nil {
					t.Error(err)
				} else {
					data, _ := ioutil.ReadAll(fd.Reader)
					t.Logf("%s %d bytes", fd.ContentType, len(data))
				}
			}
		}
	}

	if recs, err := app.GetRecords(testDataRecords.input[0].([]string), testDataRecords.input[1].(string)); err != nil {
		t.Error(err)
	} else {
		if len(recs) > 3 {
			t.Error("Too many records")
		}
	}

	if recs, err := app.GetAllRecords(testDataRecords.input[0].([]string)); err != nil {
		t.Error(err)
	} else {
		t.Log(len(recs))
	}
}

func TestGetRecordWithTotalCount(t *testing.T) {
	testDataRecords := GetTestDataGetRecordsWithTotalCount()
	app := newApp()

	if recs, totalCount, err := app.GetRecordsWithTotalCount(testDataRecords.input[0].([]string), testDataRecords.input[1].(string)); err != nil {
		t.Error(err)
	} else {
		if len(recs) > 3 {
			t.Error("Too many records")
		}
		if totalCount != "999" {
			t.Error("TotalCount incorrect", err)
		}
	}
}

func TestUpdateRecord(t *testing.T) {
	testData := GetTestDataGetRecord()
	testDataRecords := GetTestDataGetRecords()
	testDataRecordByKey := GetTestDataUpdateRecordByKey()

	app := newApp()

	rec, err := app.GetRecord(uint64(testData.input[0].(int)))
	if err != nil {
		t.Fatal(err)
	}
	rec.Fields["title"] = SingleLineTextField("new title")
	if err := app.UpdateRecord(rec, testData.input[1].(bool)); err != nil {
		t.Error("UpdateRecord failed", err)
	}

	rec.Fields[testDataRecordByKey.input[1].(string)] = SingleLineTextField(` {
		"field": "unique_key",
		"value": "unique_code"
	}`)
	if err := app.UpdateRecordByKey(rec, testData.input[1].(bool), testDataRecordByKey.input[1].(string)); err != nil {
		t.Error("UpdateRecordByKey failed", err)
	}
	recs, err := app.GetRecords(testDataRecords.input[0].([]string), testDataRecords.input[1].(string))
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
	if err := app.UpdateRecords(recs, testData.input[1].(bool)); err != nil {
		t.Error("UpdateRecords failed", err)
	}

	if err := app.UpdateRecordsByKey(recs, testDataRecordByKey.input[2].(bool), testDataRecordByKey.input[1].(string)); err != nil {
		t.Error("UpdateRecordsByKey failed", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	testData := GetTestDataDeleteRecords()
	app := newApp()
	if err := app.DeleteRecords(testData.input[0].([]uint64)); err != nil {
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
	app := newApp()
	fi, err := app.Fields()
	if err != nil {
		t.Error("Fields failed", err)
	}
	for _, f := range fi {
		t.Log(f)
	}
}

func TestApiToken(t *testing.T) {
	app := newAppWithToken()
	_, err := app.Fields()
	if err != nil {
		t.Error("Api token failed", err)
	}
}

func TestGuestSpace(t *testing.T) {
	app := newAppWithGuest()
	_, err := app.Fields()
	if err != nil {
		t.Error("GuestSpace failed", err)
	}
}

func TestGetRecordComments(t *testing.T) {
	testData := GetDataTestRecordComments()
	app := newApp()
	if rec, err := app.GetRecordComments(uint64(testData.input[0].(int)), testData.input[1].(string), uint64(testData.input[2].(int)), uint64(testData.input[3].(int))); err != nil {
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
	testData := GetDataTestDeleteRecordComment()
	appTest := newApp()
	err := appTest.DeleteComment(uint64(testData.input[0].(int)), uint64(testData.input[1].(int)))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("The comment with id =  %v has been deleted successefully!", uint64(testData.input[0].(int)))
	}
}

func TestGetProcess(t *testing.T) {
	TestData := GetTestDataProcess()
	app := newApp()
	_, err := app.GetProcess(TestData.input[0].(string))
	if err != nil {
		t.Error("TestGetProcess failed: ", err)
	}
}

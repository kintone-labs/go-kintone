// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func newApp(appId uint64) *App {
	return &App{
		Domain:   os.Getenv("KINTONE_DOMAIN"),
		User:     os.Getenv("KINTONE_USER"),
		Password: os.Getenv("KINTONE_PASSWORD"),
		AppId:    appId,
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
		if rec.Id() != "116" {
			t.Errorf("Unexpected Id: ", rec.Id())
		}
		for _, f := range rec {
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

	rec := Record{
		"title": SingleLineTextField("test!"),
		"file": FileField{
			{FileKey: fileKey},
		},
	}
	_, err = a.AddRecord(rec)
	if err != nil {
		t.Error("AddRecord failed", rec)
	}

	recs := []Record{
		{"title": SingleLineTextField("multi add 1")},
		{"title": SingleLineTextField("multi add 2")},
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

	rec := Record{
		"title": SingleLineTextField("new title"),
	}
	if err := a.UpdateRecord(4, rec); err != nil {
		t.Error("UpdateRecord failed", err)
	}

	recs := map[uint64]Record{
		1: {"title": SingleLineTextField("updated 1")},
		2: {"title": SingleLineTextField("updated 2")},
		3: {"title": SingleLineTextField("updated 3")},
	}
	if err := a.UpdateRecords(recs); err != nil {
		t.Error("UpdateRecords failed", err)
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

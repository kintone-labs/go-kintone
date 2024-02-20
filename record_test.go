// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"testing"
	"time"
)

func TestNumericId(t *testing.T) {
	if rev, err := numericId("33"); err != nil {
		t.Error(err)
	} else {
		if rev != 33 {
			t.Error("33")
		}
	}
	if rev, err := numericId("INFRA-34"); err != nil {
		t.Error(err)
	} else {
		if rev != 34 {
			t.Error("INFRA-34")
		}
	}
	if _, err := numericId("INFRA_35"); err == nil {
		t.Error("INFRA_35 must fail")
	}
}

func TestDecodeRecord(t *testing.T) {
	t.Parallel()

	j := []byte(`
{
    "record": {
        "record_id": {
            "type": "RECORD_NUMBER",
            "value": "1"
        },
        "created_time": {
            "type": "CREATED_TIME",
            "value": "2012-02-03T08:50:00Z"
        },
        "updated_time": {
            "type": "UPDATED_TIME",
            "value": "2018-10-24T08:50:00Z"
        },
        "dropdown": {
            "type": "DROP_DOWN",
            "value": "Option1"
        },
        "1line": {
            "type": "SINGLE_LINE_TEXT",
            "value": "hoge"
        },
        "2line": {
            "type": "MULTI_LINE_TEXT",
            "value": "hoge\nfuga"
        },
        "number": {
            "type": "NUMBER",
            "value": "123.456"
        },
        "check_box": {
            "type": "CHECK_BOX",
            "value": ["a", "b"]
        },
        "file": {
            "type": "FILE",
            "value": [
    {
        "contentType": "text/plain",
        "fileKey":"201202061155587E339F9067544F1A92C743460E3D12B3297",
        "name": "17to20_VerupLog (1).txt",
        "size": "12345"
    },
    {
        "contentType": "text/plain",
        "fileKey": "201202061155583C763E30196F419E83E91D2E4A03746C273",
        "name": "17to20_VerupLog.txt",
        "size": "23175"
    }
]
        },
        "date": {
            "type": "DATE",
            "value": "1974-04-04"
        },
        "time": {
            "type": "TIME",
            "value": "09:53"
        },
        "time2": {
            "type": "TIME",
            "value": "10:04:37"
        },
        "datetime": {
            "type": "DATETIME",
            "value": "2012-01-11T11:30:00Z"
        },
        "user": {
            "type": "USER_SELECT",
            "value": [
              {
                "code": "sato",
                "name": "Noboru Sato"
              }
            ]
        },
        "org": {
            "type": "ORGANIZATION_SELECT",
            "value": [
              {
                "code": "sales",
                "name": "Sales Dept"
              }
            ]
        },
        "group": {
            "type": "GROUP_SELECT",
            "value": [
              {
                "code": "managers",
                "name": "Manager Group"
              }
            ]
        },
        "$revision": {
            "type": "__REVISION__",
            "value": "7"
        }
    }
}`)
	rec, err := DecodeRecord(j)
	fields := rec.Fields
	if err != nil {
		t.Fatal(err)
	}
	if rec.Revision() != 7 {
		t.Errorf("rec.Revision() != 7 (%d)", rec.Revision())
	}
	if _, ok := fields["record_id"].(RecordNumberField); !ok {
		t.Error("Not a RecordNumberField")
	}
	if fields["record_id"] != RecordNumberField("1") {
		t.Error("record_id mismatch")
	}
	ctime, ok := fields["created_time"].(CreationTimeField)
	if !ok {
		t.Error("Not a CreationTimeField")
	}
	if time.Time(ctime).Year() != 2012 {
		t.Error("Year != 2012")
	}
	if time.Time(ctime).Month() != time.February {
		t.Error("Month != February")
	}
	if time.Time(ctime).Day() != 3 {
		t.Error("Day != 3")
	}
	if time.Time(ctime).Hour() != 8 {
		t.Error("Hour != 8")
	}
	if time.Time(ctime).Minute() != 50 {
		t.Error("Minute != 50")
	}
	if time.Time(ctime).Second() != 0 {
		t.Error("Second != 0")
	}
	mtime, ok := fields["updated_time"].(ModificationTimeField)
	if !ok {
		t.Error("Not a ModificationTimeField")
	}
	if time.Time(mtime).Year() != 2018 {
		t.Error("Year != 2018")
	}
	if time.Time(mtime).Month() != time.October {
		t.Error("Month != October")
	}
	if time.Time(mtime).Day() != 24 {
		t.Error("Day != 24")
	}
	if time.Time(mtime).Hour() != 8 {
		t.Error("Hour != 8")
	}
	if time.Time(mtime).Minute() != 50 {
		t.Error("Minute != 50")
	}
	if time.Time(mtime).Second() != 0 {
		t.Error("Second != 0")
	}
	dropdown, ok := fields["dropdown"].(SingleSelectField)
	if !ok {
		t.Error("Not a SingleSelectField")
	}
	if !dropdown.Valid {
		t.Error("Invalid dropdown")
	}
	if dropdown.String != "Option1" {
		t.Error("dropdown mismatch")
	}
	if _, ok := fields["1line"].(SingleLineTextField); !ok {
		t.Error("Not a SingleLineTextField")
	}
	if fields["1line"] != SingleLineTextField("hoge") {
		t.Error("1line mismatch")
	}
	if _, ok := fields["2line"].(MultiLineTextField); !ok {
		t.Error("Not a MultiLineTextField")
	}
	if fields["2line"] != MultiLineTextField("hoge\nfuga") {
		t.Error("2line mismatch")
	}
	num, ok := fields["number"].(DecimalField)
	if !ok {
		t.Error("Not a DecimalField")
	}
	if num != DecimalField("123.456") {
		t.Error("number mismatch")
	}
	check_box, ok := fields["check_box"].(CheckBoxField)
	if !ok {
		t.Error("Not a CheckBoxField")
	}
	if len(check_box) != 2 {
		t.Error("check_box mismatch")
	}
	file, ok := fields["file"].(FileField)
	if !ok {
		t.Error("Not a FileField")
	}
	if file[0].Size != 12345 {
		t.Error("file size mismatch")
	}
	if file[1].Name != "17to20_VerupLog.txt" {
		t.Error("file name mismatch")
	}
	date, ok := fields["date"].(DateField)
	if !ok {
		t.Error("Not a DateField")
	}
	if !date.Valid {
		t.Error("date invalid")
	}
	if date.Date.Year() != 1974 {
		t.Error("Year != 1974")
	}
	if date.Date.Month() != time.April {
		t.Error("Month != April")
	}
	if date.Date.Day() != 4 {
		t.Error("Day != 4")
	}
	time1, ok := fields["time"].(TimeField)
	if !ok {
		t.Error("Not a TimeField")
	}
	if !time1.Valid {
		t.Error("time1 invalid")
	}
	if time1.Time.Hour() != 9 {
		t.Error("Hour != 9")
	}
	if time1.Time.Minute() != 53 {
		t.Error("Minute != 53")
	}
	if time1.Time.Second() != 0 {
		t.Error("Second != 0")
	}
	time2, ok := fields["time2"].(TimeField)
	if !ok {
		t.Error("Not a TimeField")
	}
	if !time2.Valid {
		t.Error("time2 invalid")
	}
	if time2.Time.Hour() != 10 {
		t.Error("Hour != 10")
	}
	if time2.Time.Minute() != 4 {
		t.Error("Minute != 4")
	}
	if time2.Time.Second() != 37 {
		t.Error("Second != 37")
	}
	dt, ok := fields["datetime"].(DateTimeField)
	if !ok {
		t.Error("Not a DateTimeField")
	}
	if dt.Time.Hour() != 11 {
		t.Error("Hour != 11")
	}
	if dt.Time.Minute() != 30 {
		t.Error("Minute != 30")
	}
	user, ok := fields["user"].(UserField)
	if !ok {
		t.Error("Not a UserField")
	}
	if user[0].Code != "sato" {
		t.Error("user code mismatch")
	}
	if user[0].Name != "Noboru Sato" {
		t.Error("user name mismatch")
	}
	org, ok := fields["org"].(OrganizationField)
	if !ok {
		t.Error("Not a OrganizationField")
	}
	if org[0].Code != "sales" {
		t.Error("organization code mismatch")
	}
	if org[0].Name != "Sales Dept" {
		t.Error("organization name mismatch")
	}
	group, ok := fields["group"].(GroupField)
	if !ok {
		t.Error("Not a GroupField")
	}
	if group[0].Code != "managers" {
		t.Error("group code mismatch")
	}
	if group[0].Name != "Manager Group" {
		t.Error("group name mismatch")
	}
}

func TestDecodeRecords(t *testing.T) {
	t.Parallel()

	j := []byte(`
{
    "records": [
        {
            "record_id": {
                "type": "RECORD_NUMBER",
                "value": "1"
            },
            "created_time": {
                "type": "CREATED_TIME",
                "value": "2012-02-03T08:50:00Z"
            },
            "updated_time": {
                "type": "UPDATED_TIME",
                "value": "2018-10-24T08:50:00Z"
            },
            "dropdown": {
                "type": "DROP_DOWN",
                "value": null
            }
        },
        {
            "record_id": {
                "type": "RECORD_NUMBER",
                "value": "2"
            },
            "created_time": {
                "type": "CREATED_TIME",
                "value": "2012-02-03T09:22:00Z"
            },
            "updated_time": {
                "type": "UPDATED_TIME",
                "value": "2018-10-24T09:22:00Z"
            },
            "dropdown": {
                "type": "DROP_DOWN",
                "value": null
            }
        }
    ]
}`)
	rec, err := DecodeRecords(j)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec) != 2 {
		t.Error("length mismatch")
	}
	if _, ok := rec[0].Fields["record_id"]; !ok {
		t.Error("record_id must exist")
	}
	dropdown, ok := rec[0].Fields["dropdown"]
	if !ok {
		t.Error("null dropdown field must exist")
	}
	if dropdown.(SingleSelectField).Valid {
		t.Error("dropdown must be invalid")
	}
}

func TestDecodeRecordsWithTotalCount(t *testing.T) {
	b := []byte(`
	{
		"records": [
			{
				"record_id": {
					"type": "RECORD_NUMBER",
					"value": "1"
				},
				"created_time": {
					"type": "CREATED_TIME",
					"value": "2012-02-03T08:50:00Z"
				},
				"updated_time": {
					"type": "UPDATED_TIME",
					"value": "2018-10-24T08:50:00Z"
				},
				"dropdown": {
					"type": "DROP_DOWN",
					"value": null
				}
			},
			{
				"record_id": {
					"type": "RECORD_NUMBER",
					"value": "2"
				},
				"created_time": {
					"type": "CREATED_TIME",
					"value": "2012-02-03T09:22:00Z"
				},
				"updated_time": {
					"type": "UPDATED_TIME",
					"value": "2018-10-24T09:22:00Z"
				},
				"dropdown": {
					"type": "DROP_DOWN",
					"value": null
				}
			}
		],
		"totalCount": "9999"
	}`)

	rec, totalCount, err := DecodeRecordsWithTotalCount(b)
	if err != nil {
		t.Fatal(err)
	}
	if totalCount != "9999" {
		t.Error("totalCount is incorrect")
	}
	if len(rec) != 2 {
		t.Error("length mismatch")
	}
	if _, ok := rec[0].Fields["record_id"]; !ok {
		t.Error("record_id must exist")
	}
	dropdown, ok := rec[0].Fields["dropdown"]
	if !ok {
		t.Error("null dropdown field must exist")
	}
	if dropdown.(SingleSelectField).Valid {
		t.Error("dropdown must be invalid")
	}
}

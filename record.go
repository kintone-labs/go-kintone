// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Record represens a record in an application.
//
// In fact, Record is a mapping between field IDs and fields.
// Although field types are shown as interface{}, they are guaranteed
// to be one of a *Field type in this package.
type Record map[string]interface{}

// Id returns the record number.
//
// A record number is unique within an application.
func (rec Record) Id() string {
	for _, field := range map[string]interface{}(rec) {
		if id, ok := field.(RecordNumberField); ok {
			return string(id)
		}
	}
	panic("No record number field")
}

// Assert string list.
func stringList(l []interface{}) []string {
	sl := make([]string, len(l))
	for i, v := range l {
		sl[i] = v.(string)
	}
	return sl
}

// Convert user list.
func userList(l []interface{}) ([]User, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	var ul []User
	err = json.Unmarshal(b, &ul)
	if err != nil {
		return nil, err
	}
	return ul, nil
}

type recordData map[string]struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func decodeRecordData(data recordData) (Record, error) {
	rec := make(map[string]interface{})
	for key, v := range data {
		switch v.Type {
		case FT_SINGLE_LINE_TEXT:
			rec[key] = SingleLineTextField(v.Value.(string))
		case FT_MULTI_LINE_TEXT:
			rec[key] = MultiLineTextField(v.Value.(string))
		case FT_RICH_TEXT:
			rec[key] = RichTextField(v.Value.(string))
		case FT_DECIMAL:
			rec[key] = DecimalField(v.Value.(string))
		case FT_CALC:
			rec[key] = CalcField(v.Value.(string))
		case FT_CHECK_BOX:
			rec[key] = CheckBoxField(stringList(v.Value.([]interface{})))
		case FT_RADIO:
			rec[key] = RadioButtonField(v.Value.(string))
		case FT_SINGLE_SELECT:
			if v.Value == nil {
				rec[key] = SingleSelectField{Valid: false}
			} else {
				rec[key] = SingleSelectField{v.Value.(string), true}
			}
		case FT_MULTI_SELECT:
			rec[key] = MultiSelectField(stringList(v.Value.([]interface{})))
		case FT_FILE:
			b1, err := json.Marshal(v.Value)
			if err != nil {
				return nil, err
			}
			var fl []File
			err = json.Unmarshal(b1, &fl)
			if err != nil {
				return nil, err
			}
			rec[key] = FileField(fl)
		case FT_LINK:
			rec[key] = LinkField(v.Value.(string))
		case FT_DATE:
			if v.Value == nil {
				rec[key] = DateField{Valid: false}
			} else {
				d, err := time.Parse("2006-01-02", v.Value.(string))
				if err != nil {
					return nil, fmt.Errorf("Invalid date: %v", v.Value)
				}
				rec[key] = DateField{d, true}
			}
		case FT_TIME:
			if v.Value == nil {
				rec[key] = TimeField{Valid: false}
			} else {
				t, err := time.Parse("15:04", v.Value.(string))
				if err != nil {
					t, err = time.Parse("15:04:05", v.Value.(string))
					if err != nil {
						return nil, fmt.Errorf("Invalid time: %v", v.Value)
					}
				}
				rec[key] = TimeField{t, true}
			}
		case FT_DATETIME:
			if s, ok := v.Value.(string); ok {
				dt, err := time.Parse(time.RFC3339, s)
				if err != nil {
					return nil, fmt.Errorf("Invalid datetime: %v", v.Value)
				}
				rec[key] = DateTimeField(dt)
			}
		case FT_USER:
			ul, err := userList(v.Value.([]interface{}))
			if err != nil {
				return nil, err
			}
			rec[key] = UserField(ul)
		case FT_CATEGORY:
			rec[key] = CategoryField(stringList(v.Value.([]interface{})))
		case FT_STATUS:
			rec[key] = StatusField(v.Value.(string))
		case FT_ASSIGNEE:
			al, err := userList(v.Value.([]interface{}))
			if err != nil {
				return nil, err
			}
			rec[key] = AssigneeField(al)
		case FT_ID:
			rec[key] = RecordNumberField(v.Value.(string))
		case FT_CREATOR:
			creator := v.Value.(map[string]interface{})
			rec[key] = CreatorField{
				creator["code"].(string),
				creator["name"].(string),
			}
		case FT_CTIME:
			var ctime time.Time
			if ctime.UnmarshalText([]byte(v.Value.(string))) != nil {
				return nil, fmt.Errorf("Invalid datetime: %v", v.Value)
			}
			rec[key] = CreationTimeField(ctime)
		case FT_MODIFIER:
			modifier := v.Value.(map[string]interface{})
			rec[key] = ModifierField{
				modifier["code"].(string),
				modifier["name"].(string),
			}
		case FT_MTIME:
			var mtime time.Time
			if mtime.UnmarshalText([]byte(v.Value.(string))) != nil {
				return nil, fmt.Errorf("Invalid datetime: %v", v.Value)
			}
			rec[key] = CreationTimeField(mtime)
		case FT_SUBTABLE:
			b2, err := json.Marshal(v.Value)
			if err != nil {
				return nil, err
			}
			var stl []SubTableEntry
			err = json.Unmarshal(b2, &stl)
			if err != nil {
				return nil, err
			}
			rec[key] = SubTableField(stl)
		default:
			return nil, fmt.Errorf("Invalid type: %v", v.Type)
		}
	}
	return rec, nil
}

// DecodeRecords decodes JSON response for multi-get API.
func DecodeRecords(b []byte) ([]Record, error) {
	var t struct {
		Records []recordData `json:"records"`
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}
	rec_list := make([]Record, len(t.Records))
	for i, rd := range t.Records {
		r, err := decodeRecordData(rd)
		if err != nil {
			return nil, err
		}
		rec_list[i] = r
	}
	return rec_list, nil
}

// DecodeRecord decodes JSON response for single-get API.
func DecodeRecord(b []byte) (Record, error) {
	var t struct {
		RecordData recordData `json:"record"`
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}
	rec, err := decodeRecordData(t.RecordData)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

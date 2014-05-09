// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"encoding/json"
	"strconv"
	"time"
)

// Field type identifiers.
const (
	FT_SINGLE_LINE_TEXT = "SINGLE_LINE_TEXT"
	FT_MULTI_LINE_TEXT  = "MULTI_LINE_TEXT"
	FT_RICH_TEXT        = "RICH_TEXT"
	FT_DECIMAL          = "NUMBER"
	FT_CALC             = "CALC"
	FT_CHECK_BOX        = "CHECK_BOX"
	FT_RADIO            = "RADIO_BUTTON"
	FT_SINGLE_SELECT    = "DROP_DOWN"
	FT_MULTI_SELECT     = "MULTI_SELECT"
	FT_FILE             = "FILE"
	FT_LINK             = "LINK"
	FT_DATE             = "DATE"
	FT_TIME             = "TIME"
	FT_DATETIME         = "DATETIME"
	FT_USER             = "USER_SELECT"
	FT_CATEGORY         = "CATEGORY"
	FT_STATUS           = "STATUS"
	FT_ASSIGNEE         = "STATUS_ASSIGNEE"
	FT_RECNUM           = "RECORD_NUMBER"
	FT_CREATOR          = "CREATOR"
	FT_CTIME            = "CREATED_TIME"
	FT_MODIFIER         = "MODIFIER"
	FT_MTIME            = "UPDATED_TIME"
	FT_SUBTABLE         = "SUBTABLE"
	FT_ID               = "__ID__"
	FT_REVISION         = "__REVISION__"
)

// SingleLineTextField is a field type for single-line texts.
type SingleLineTextField string

func (f SingleLineTextField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_SINGLE_LINE_TEXT,
		"value": string(f),
	})
}

// MultiLineTextField is a field type for multi-line texts.
type MultiLineTextField string

func (f MultiLineTextField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_MULTI_LINE_TEXT,
		"value": string(f),
	})
}

// RichTextField is a field type for HTML rich texts.
type RichTextField string

func (f RichTextField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_RICH_TEXT,
		"value": string(f),
	})
}

// DecimalField is a field type for decimal numbers.
type DecimalField string

func (f DecimalField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_DECIMAL,
		"value": string(f),
	})
}

// CalcField is a field type for auto-calculated values.
type CalcField string

func (f CalcField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_CALC,
		"value": string(f),
	})
}

// CheckBoxField is a field type for selected values in a check-box.
type CheckBoxField []string

func (f CheckBoxField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_CHECK_BOX,
		"value": []string(f),
	})
}

// RadioButtonField is a field type for the selected value by a radio-button.
type RadioButtonField string

func (f RadioButtonField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_RADIO,
		"value": string(f),
	})
}

// SingleSelectField is a field type for the selected value in a selection box.
type SingleSelectField struct {
	String string // Selected value.
	Valid  bool   // If not selected, false.
}

func (f SingleSelectField) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return json.Marshal(map[string]interface{}{
			"type":  FT_SINGLE_SELECT,
			"value": f.String,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"type":  FT_SINGLE_SELECT,
			"value": nil,
		})
	}
}

// MultiSelectField is a field type for selected values in a selection box.
type MultiSelectField []string

func (f MultiSelectField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_MULTI_SELECT,
		"value": []string(f),
	})
}

// File is a struct for an uploaded file.
type File struct {
	ContentType string `json:"contentType"` // MIME type of the file
	FileKey     string `json:"fileKey"`     // BLOB ID of the file
	Name        string `json:"name"`        // File name
	Size        uint64 `json:"size,string"` // The file size
}

func (f *File) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		map[string]interface{}{
			"contentType": f.ContentType,
			"fileKey":     f.FileKey,
			"name":        f.Name,
			"size":        strconv.FormatUint(f.Size, 10),
		})
}

// FileField is a field type for uploaded files.
type FileField []File

func (f FileField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_FILE,
		"value": []File(f),
	})
}

// LinkField is a field type for hyper-links.
type LinkField string

func (f LinkField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_LINK,
		"value": string(f),
	})
}

// DateField is a field type for dates.
type DateField struct {
	Date  time.Time // stores date information.
	Valid bool      // false when not set.
}

// NewDateField returns an instance of DateField.
func NewDateField(year int, month time.Month, day int) DateField {
	return DateField{
		time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		true,
	}
}

func (f DateField) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return json.Marshal(map[string]interface{}{
			"type":  FT_DATE,
			"value": f.Date.Format("2006-01-02"),
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"type":  FT_DATE,
			"value": nil,
		})
	}
}

// TimeField is a field type for times.
type TimeField struct {
	Time  time.Time // stores time information.
	Valid bool      // false when not set.
}

// NewTimeField returns an instance of TimeField.
func NewTimeField(hour, min int) TimeField {
	return TimeField{
		time.Date(1, time.January, 1, hour, min, 0, 0, time.UTC),
		true,
	}
}

func (f TimeField) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return json.Marshal(map[string]interface{}{
			"type":  FT_TIME,
			"value": f.Time.Format("15:04:05"),
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"type":  FT_TIME,
			"value": nil,
		})
	}
}

// DateTimeField is a field type for date & time.
type DateTimeField time.Time

// NewDateTimeField returns an instance of DateTimeField.
func NewDateTimeField(year int, month time.Month, day, hour, min int) DateTimeField {
	return DateTimeField(
		time.Date(year, month, day, hour, min, 0, 0, time.UTC))
}

func (f DateTimeField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_DATETIME,
		"value": time.Time(f).Format(time.RFC3339),
	})
}

// User represents a user entry.
type User struct {
	Code string `json:"code"` // A unique identifer of the user.
	Name string `json:"name"` // The user name.
}

// UserField is a field type for user entries.
type UserField []User

func (f UserField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_USER,
		"value": []User(f),
	})
}

// CategoryField is a list of category names.
type CategoryField []string

func (f CategoryField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_CATEGORY,
		"value": []string(f),
	})
}

// StatusField is a string label of a record status.
type StatusField string

func (f StatusField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_STATUS,
		"value": string(f),
	})
}

// AssigneeField is a list of user entries who are assigned to a record.
type AssigneeField []User

func (f AssigneeField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_ASSIGNEE,
		"value": []User(f),
	})
}

// RecordNumberField is a record number.
type RecordNumberField string

func (f RecordNumberField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_RECNUM,
		"value": string(f),
	})
}

// CreatorField is a user who created a record.
type CreatorField User

func (f CreatorField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_CREATOR,
		"value": User(f),
	})
}

// CreationTimeField is the time when a record is created.
type CreationTimeField time.Time

func (t CreationTimeField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_CTIME,
		"value": time.Time(t).Format(time.RFC3339),
	})
}

// ModifierField is a user who modified a record last.
type ModifierField User

func (f ModifierField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_MODIFIER,
		"value": User(f),
	})
}

// ModificationTimeField is the time when a record is last modified.
type ModificationTimeField time.Time

func (t ModificationTimeField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_MTIME,
		"value": time.Time(t).Format(time.RFC3339),
	})
}

// SubTableEntry is a type for an entry in a subtable.
type SubTableEntry struct {
	Id    string                 `json:"id"`    // The entry ID
	Value map[string]interface{} `json:"value"` // Subtable data fields.
}

// SubTableField is a list of subtable entries.
type SubTableField []SubTableEntry

func (f SubTableField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  FT_SUBTABLE,
		"value": []SubTableEntry(f),
	})
}

// IsBuiltinField returns true if the field is a built-in field.
func IsBuiltinField(o interface{}) bool {
	switch o.(type) {
	case CalcField:
		return true
	case CategoryField:
		return true
	case StatusField:
		return true
	case AssigneeField:
		return true
	case RecordNumberField:
		return true
	case CreatorField:
		return true
	case CreationTimeField:
		return true
	case ModifierField:
		return true
	case ModificationTimeField:
		return true
	}
	return false
}

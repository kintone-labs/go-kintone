package kintone

import (
	"encoding/json"
)

//Object Cursor structure
type Cursor struct {
	Id         string `json:"id"`
	TotalCount string `json:"totalCount"`
}
type RecordCursor struct {
	Records []*Record `json:"records"`
	Next    bool      `json:"next"`
}

//decodeCursor decodes JSON response for cursor api
func decodeCursor(b []byte) (c *Cursor, err error) {
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func decodeRecordCursor(b []byte) (rc *RecordCursor, err error) {
	var t struct {
		next bool
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	listRecord, err := DecodeRecords(b)
	if err != nil {
		return nil, err
	}
	records := &RecordCursor{Records: listRecord, Next: t.next}
	return records, nil
}

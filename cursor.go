package kintone

import (
	"encoding/json"
)

//Object Cursor structure
type Cursor struct {
	Id         string `json:"id"`
	TotalCount string `json:"totalCount"`
}
type GetRecordsCursorResponse struct {
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
func DecodeGetRecordsCursorResponse(b []byte) (rc *GetRecordsCursorResponse, err error) {
	var t struct {
		Next bool `json:"next"`
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	listRecord, err := DecodeRecords(b)
	if err != nil {
		return nil, err
	}
	getRecordsCursorResponse := &GetRecordsCursorResponse{Records: listRecord, Next: t.Next}
	return getRecordsCursorResponse, nil
}

package kintone

import (
	"encoding/json"
)

// Cursor structure
type Cursor struct {
	Id          string `json:"id"`
	TotalCount  uint64 `json:"totalCount,string"`
}

// DecodeCursor decodes JSON response for cursor api
func DecodeCursor(b []byte) (cursor *Cursor, err error) {
	err = json.Unmarshal(b, &cursor)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

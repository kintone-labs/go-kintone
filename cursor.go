package kintone

import (
	"encoding/json"
)

//Object Cursor structure
type Cursor struct {
	Id         string `json:"id"`
	TotalCount string `json:"totalCount"`
}

//decodeCursor decodes JSON response for cursor api
func decodeCursor(b []byte) (c *Cursor, err error) {
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

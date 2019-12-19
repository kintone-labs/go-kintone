package kintone

import (
	"encoding/json"
)

type Cursor struct {
	Id         string `json:"id"`
	TotalCount string `json:"totalCount"`
}

func decodeCursor(b []byte) (c *Cursor, err error) {
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

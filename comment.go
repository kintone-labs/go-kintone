package kintone

import (
	"encoding/json"
	"errors"
)

type Comment struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	CreatedAt string `json:"createdAt"`
	Creator map[string]interface{} `json:"creator"`
	Mentions []interface{} `json:"mentions"`
}

// DecodeRecordComments decodes JSON response for comment api
func DecodeRecordComments(b []byte) ([]Comment, error) {
	var t struct {
		MyComments []Comment `json:"comments"`
		Older bool  `json:"older"`
		Newer bool  `json:"newer"`
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}
	return t.MyComments, nil
}

package kintone

import (
	"testing"
)

func TestDecodeCursor(t *testing.T) {
	data := []byte(`{"id":"aaaaaaaaaaaaaaaaaa","totalCount":"null"}`)
	_, err := decodeCursor(data)
	if err != nil {
		t.Errorf("TestDecodeCursor is failed: %v", err)
	}
}

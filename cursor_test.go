package kintone

import (
	"fmt"
	"testing"
)

func TestDecodeCursor(t *testing.T) {
	data := []byte(`{"id":"aaaaaaaaaaaaaaaaaa","totalCount":"null"}`)
	cursor, err := decodeCursor(data)
	if err != nil {
		t.Errorf("TestDecodeCursor is failed: %v", err)

	}
	fmt.Println(cursor)
}

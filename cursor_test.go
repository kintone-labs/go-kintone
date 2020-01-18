package kintone

import "testing"

func TestDecodeCursor(t *testing.T) {
	j := []byte(`
{"id":"9a9716fe-1394-4677-a1c7-2199a5d28215", "totalCount": "123456"}
`)

	cursor, err := DecodeCursor(j)
	if err != nil {
		t.Fatal(err)
	}
	if cursor.Id != "9a9716fe-1394-4677-a1c7-2199a5d28215" {
		t.Errorf("cursor id mismatch. actual %v", cursor.Id)
	}
	if cursor.TotalCount != 123456 {
		t.Errorf("cursor total count mismatch. actual %v", cursor.TotalCount)
	}
}

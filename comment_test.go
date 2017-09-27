package kintone

import "testing"

func TestDecodeRecordComments(t *testing.T) {
	j := []byte(`
{"comments":[{"id":"2","text":"ほげほげ","createdAt":"2016-11-07T19:53:32Z","creator":{"code":"xxx.tat","name":"さんぷる"},"mentions":[{"code":"wolf","type":"USER"}]},{"id":"1","text":"ふがふが","createdAt":"2016-11-07T19:53:27Z","creator":{"code":"xxx.tat","name":"さんぷる"},"mentions":[]}],"older":false,"newer":false}
`)

	rec, err := DecodeRecordComments(j)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec) != 2 {
		t.Fatal("invalud record count!")
	}
	if rec[0].Id != "2" {
		t.Errorf("comment id mismatch. actual %v", rec[0].Id)
	}
	if rec[0].Text != "ほげほげ" {
		t.Errorf("comment text mismatch. actual %v", rec[0].Text)
	}

	if rec[0].CreatedAt != "2016-11-07T19:53:32Z" {
		t.Errorf("comment createdat mismatch. actual %v", rec[0].CreatedAt)
	}
	if rec[0].Creator.Code != "xxx.tat" {
		t.Errorf("comment creator mismatch. actual %v", rec[0].Creator.Code)
	}
	if rec[0].Mentions[0].Code != "wolf" {
		t.Errorf("comment mention-code mismatch. actual %v", rec[0].Mentions[0].Code)
	}
}

/*
Package kintone provides interfaces for kintone REST API.

See https://developers.cybozu.com/ja/kintone-api/app-api.html for API specs.

To retrieve 3 records from a kintone app (id=25):

	import (
		"github.com/cybozu/go-kintone"
		"log"
	)
	...
	app := kintone.App{
		"example.cybozu.com",
		"user1",
		"password",
		25,
	}
	records, err := app.GetRecords(nil, "limit 3")
	if err != nil {
		log.Fatal(err)
	}
	// use records
*/
package kintone

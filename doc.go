/*
Package kintone provides interfaces for kintone REST API.

See http://developers.kintone.com/ for API specs.

To retrieve 3 records from a kintone app (id=25):

	import (
		"log"

		"github.com/kintone/go-kintone"
	)
	...
	app := &kintone.App{
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

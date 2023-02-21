package kintone

import (
	"bytes"
)

type TestData struct {
	input  []interface{}
	output string
}

func GetTestDataProcess() *TestData {
	return &TestData{
		input: []interface{}{"en"},
		output: `
		{
			"enable":true,
			"states":{
				"Not started":{
					"name":"Not started",
					"index":"0",
					"assignee":{
						"type":"ONE",
						"entities":[]
					}
				},
				"In progress":{
					"name":"In progress",
					"index":"1",
					"assignee":{
						"type":"ALL",
						"entities":[
							{
								"entity":{
									"type":"USER",
									"code":"user1"
								},
								"includeSubs":false
							},
							{
								"entity":{
									"type":"FIELD_ENTITY",
									"code":"creator"
								},
								"includeSubs":false
							},
							{
								"entity":{
									"type":"CUSTOM_FIELD",
									"code":"Boss"
								},
								"includeSubs":false
							}
						]
					}
				},
				"Completed":{
					"name":"Completed",
					"index":"2",
					"assignee":{
						"type":"ONE",
						"entities":[]
					}
				}
			},
			"actions":[
				{
					"name":"Start",
					"from":"Not started",
					"to":"In progress",
					"filterCond":"Record_number = \"1\""
				},
				{
					"name":"Complete",
					"from":"In progress",
					"to":"Completed",
					"filterCond":""
				}
			],
			"revision":"3"
		}`,
	}
}

func GetTestDataDeleteRecords() *TestData {
	return &TestData{
		input:  []interface{}{[]uint64{6, 7}},
		output: `{}`,
	}
}

func GetTestDataGetRecord() *TestData {
	return &TestData{
		input: []interface{}{1, true},
		output: `
		{
			"record":{
				"Updated_by":{
					"type":"MODIFIER",
					"value":{
						"code":"Administrator",
						"name":"Administrator"
					},
					"key":"hehehe"
				},
				"$id":{
					"type":"__ID__",
					"value":"1"
				}
			}
		}`,
	}
}

func GetTestDataGetRecords() *TestData {
	return &TestData{
		input: []interface{}{
			[]string{},
			"limit 3 offset 3",
		},
		output: `
		{
			"records":[
				{
					"Created_datetime":{
						"type":"CREATED_TIME",
						"value":"2019-03-11T04:50:00Z"
					},
					"Created_by":{
						"type":"CREATOR",
						"value":{
							"code":"Administrator",
							"name":"Administrator"
						}
					},
					"$id":{
						"type":"__ID__",
						"value":"1"
					}
				},
				{
					"Created_datetime":{
						"type":"CREATED_TIME",
						"value":"2019-03-11T06:42:00Z"
					},
					"Created_by":{
						"type":"CREATOR",
						"value":{
							"code":"Administrator",
							"name":"Administrator"
						}
					},
					"$id":{
						"type":"__ID__",
						"value":"2"
					}
				}
			],
			"totalCount":null
		}`,
	}
}

func GetTestDataGetRecordsWithTotalCount() *TestData {
	return &TestData{
		input: []interface{}{
			[]string{},
			"limit 3 offset 3",
		},
		output: `
		{
			"records":[
				{
					"Created_datetime":{
						"type":"CREATED_TIME",
						"value":"2019-03-11T04:50:00Z"
					},
					"Created_by":{
						"type":"CREATOR",
						"value":{
							"code":"Administrator",
							"name":"Administrator"
						}
					},
					"$id":{
						"type":"__ID__",
						"value":"1"
					}
				},
				{
					"Created_datetime":{
						"type":"CREATED_TIME",
						"value":"2019-03-11T06:42:00Z"
					},
					"Created_by":{
						"type":"CREATOR",
						"value":{
							"code":"Administrator",
							"name":"Administrator"
						}
					},
					"$id":{
						"type":"__ID__",
						"value":"2"
					}
				}
			],
			"totalCount": "999"
		}`,
	}
}

func GetDataTestUploadFile() *TestData {
	return &TestData{
		output: `
		{
			"app":3,
			"id":6,
			"record":{
				"attached_file":{
					"value":[
						{
							"fileKey":" c15b3870-7505-4ab6-9d8d-b9bdbc74f5d6"
						}
					]
				}
			}
		}`,
	}
}

func GetDataTestRecordComments() *TestData {
	return &TestData{
		input: []interface{}{1, "asc", 0, 10},
		output: `
		{
			"comments":[
				{
					"id":"3",
					"text":"user14 Thank you! Looks great.",
					"createdAt":"2016-05-09T18:29:05Z",
					"creator":{
						"code":"user13",
						"name":"user13"
					},
					"mentions":[
						{
							"code":"user14",
							"type":"USER"
						}
					]
				},
				{
					"id":"2",
					"text":"user13 Global Sales APAC Taskforce \nHere is today's report.",
					"createdAt":"2016-05-09T18:27:54Z",
					"creator":{
						"code":"user14",
						"name":"user14"
					},
					"mentions":[
						{
							"code":"user13",
							"type":"USER"
						},
						{
							"code":"Global Sales_1BNZeQ",
							"type":"ORGANIZATION"
						},
						{
							"code":"APAC Taskforce_DJrvzu",
							"type":"GROUP"
						}
					]
				}
			],
			"older":false,
			"newer":false
		}`,
	}
}

func GetDataTestForm() *TestData {
	return &TestData{
		output: `
		{
			"properties":[
				{
					"code":"string_1",
					"defaultValue":"",
					"expression":"",
					"hideExpression":"false",
					"maxLength":"64",
					"minLength":null,
					"label":"string_1",
					"noLabel":"false",
					"required":"true",
					"type":"SINGLE_LINE_TEXT",
					"unique":"true"
				},
				{
					"code":"number_1",
					"defaultValue":"12345",
					"digit":"true",
					"displayScale":"4",
					"expression":"",
					"maxValue":null,
					"minValue":null,
					"label":"number_1",
					"noLabel":"true",
					"required":"false",
					"type":"NUMBER",
					"unique":"false"
				},
				{
					"code":"checkbox_1",
					"defaultValue":[
						"sample1",
						"sample3"
					],
					"label":"checkbox_1",
					"noLabel":"false",
					"options":[
						"sample1",
						"sample2",
						"sample3"
					],
					"required":"false",
					"type":"CHECK_BOX"
				}
			]
		}`,
	}
}

func GetDataTestDeleteRecordComment() *TestData {
	return &TestData{
		input:  []interface{}{3, 14},
		output: `{}`,
	}
}

func GetTestDataAddRecord() *TestData {
	return &TestData{
		output: `{
			"id": "1",
			"revision": "1"
		}`,
	}
}

func GetTestDataAddRecords() *TestData {
	return &TestData{
		output: `
		{
			"ids": ["77","78"],
			"revisions": ["1","1"]
		}`,
	}
}

func GetDataTestAddRecord() *TestData {
	return &TestData{
		input: []interface{}{
			"ほげ春巻.txta",
			bytes.NewReader([]byte(`abc
			<a href="https://www.cybozu.com/">hoge</a>
			`)),
			"text/html",
		},
		output: `
		{
		  "id": "1",
		  "revision": "1"
		}`,
	}
}

func getDataTestCreateCursor() *TestData {
	return &TestData{
		output: `
		{
		  "id": "9a9716fe-1394-4677-a1c7-2199a5d28215",
		  "totalCount": 123456
		}`,
	}
}

func GetDataTestGetRecordsByCursor() *TestData {
	return &TestData{
		input: []interface{}{"9a9716fe-1394-4677-a1c7-2199a5d28215"},
		output: `
		{
			"records":[
				{
					"$id":{
						"type":"__ID__",
						"value":"1"
					},
					"Created_by":{
						"type":"CREATOR",
						"value":{
							"code":"Administrator",
							"name":"Administrator"
						}
					},
					"Created_datetime":{
						"type":"CREATED_TIME",
						"value":"2019-05-23T04:50:00Z"
					}
				}
			],
			"next":false
		}`,
	}
}

func GetTestDataDeleteCursor() *TestData {
	return &TestData{
		input:  []interface{}{"9a9716fe-1394-4677-a1c7-2199a5d28215"},
		output: `{}`,
	}
}

func GetTestDataCreateCursor() *TestData {
	return &TestData{
		input:  []interface{}{[]string{"$id", "date"}, "", 100},
		output: `{"id":"9a9716fe-1394-4677-a1c7-2199a5d28215"}`,
	}
}

func GetTestDataAddRecordComment() *TestData {
	return &TestData{
		input:  []interface{}{2},
		output: `{"id": "4"}`,
	}
}

func GetTestDataUpdateRecordByKey() *TestData {
	return &TestData{
		input: []interface{}{2, "key", true},
		output: `
		{
			"app":1,
			"records":[
				{
					"updateKey":{
						"field":"unique_key",
						"value":"CODE123"
					},
					"record":{
						"Text":{
							"value":"Silver plates"
						}
					}
				},
				{
					"updateKey":{
						"field":"unique_key",
						"value":"CODE456"
					},
					"record":{
						"Text":{
							"value":"The quick brown fox."
						}
					}
				}
			]
		}`,
	}
}

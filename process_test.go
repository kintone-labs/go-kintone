// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"testing"
)

func TestDecodeProcess(t *testing.T) {
	t.Parallel()
	j := []byte(`
	{
		"enable": true,
		"states": {
			"Not started": {
				"name": "Not started",
				"index": "0",
				"assignee": {
					"type": "ONE",
					"entities": [
					]
				}
			},
			"In progress": {
				"name": "In progress",
				"index": "1",
				"assignee": {
					"type": "ALL",
					"entities": [
						{
							"entity": {
								"type": "USER",
								"code": "user1"
							},
							"includeSubs": false
						},
						{
							"entity": {
								"type": "FIELD_ENTITY",
								"code": "creator"
							},
							"includeSubs": false
						},
						{
							"entity": {
								"type": "CUSTOM_FIELD",
								"code": "Boss"
							},
							"includeSubs": false
						}
					]
				}
			},
			"Completed": {
				"name": "Completed",
				"index": "2",
				"assignee": {
					"type": "ONE",
					"entities": [
					]
				}
			}
		},
		"actions": [
			{
				"name": "Start",
				"from": "Not started",
				"to": "In progress",
				"filterCond": "Record_number = \"1\""
			},
			{
				"name": "Complete",
				"from": "In progress",
				"to": "Completed",
				"filterCond": ""
			}
		],
		"revision": "3"
	}`)
	process, err := DecodeProcess(j)
	if err != nil {
		t.Fatal(err)
	}
	if !process.Enable {
		t.Errorf("Expected enabled, got %v", process.Enable)
	}
	if len(process.States) != 3 {
		t.Errorf("Expected 3 states, got %v", len(process.States))
	}
	if len(process.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %v", len(process.Actions))
	}
	if process.Revision != "3" {
		t.Errorf("Expected revision to be 3, got %v", process.Revision)
	}
	notStartedProcess := process.States["Not started"]
	if notStartedProcess.Name != "Not started" {
		t.Errorf("Expected status name to be \"Not started\", got \"%v\"", notStartedProcess.Name)
	}
	if notStartedProcess.Index != "0" {
		t.Errorf("Expected status index to be 0, got %v", notStartedProcess.Index)
	}
	if notStartedProcess.Assignee.Type != "ONE" {
		t.Errorf("Expected assignee type to be \"ONE\", got \"%v\"", notStartedProcess.Assignee.Type)
	}
	inProgressProcess := process.States["In progress"]
	if len(inProgressProcess.Assignee.Entities) != 3 {
		t.Errorf("Expected 0 assignees, got %v", len(inProgressProcess.Assignee.Entities))
	}
	if inProgressProcess.Assignee.Entities[0].Entity.Code != "user1" {
		t.Errorf("Expected entity to be named \"user1\", got \"%v\"", inProgressProcess.Assignee.Entities[0].Entity.Code)
	}
	if inProgressProcess.Assignee.Entities[0].IncludeSubs {
		t.Errorf("Expected entity to not include subs, got %v", inProgressProcess.Assignee.Entities[0].IncludeSubs)
	}
	completeAction := process.Actions[1]
	if completeAction.Name != "Complete" {
		t.Errorf("Expected action name to be \"Complete\", got \"%v\"", completeAction.Name)
	}
	if completeAction.From != "In progress" {
		t.Errorf("Expected previous status to be \"In progress\", got \"%v\"", completeAction.From)
	}
	if completeAction.To != "Completed" {
		t.Errorf("Expected next status to be \"Completed\", got \"%v\"", completeAction.To)
	}
}

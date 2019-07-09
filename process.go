// (C) 2014 Cybozu.  All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package kintone

import (
	"encoding/json"
)

// Process represents the process management settings for an application
type Process struct {
	Enable   bool                       `json:"enable"`
	States   map[string](*ProcessState) `json:"states"`
	Actions  []*ProcessAction           `json:"actions"`
	Revision string                     `json:"revision"`
}

// ProcessState represents a process management status
type ProcessState struct {
	Name     string           `json:"name"`
	Index    string           `json:"index"`
	Assignee *ProcessAssignee `json:"assignee"`
}

// ProcessAction representes a process management action
type ProcessAction struct {
	Name       string `json:"name"`
	From       string `json:"from"`
	To         string `json:"to"`
	FilterCond string `json:"filterCond"`
}

// ProcessAssignee represents a ProcessState assignee
type ProcessAssignee struct {
	Type     string           `json:type`
	Entities []*ProcessEntity `json:entities`
}

// ProcessEntity represents a process assignee entity
type ProcessEntity struct {
	Entity      *Entity `json:entity`
	IncludeSubs bool    `json:includeSubs`
}

// Entity is the concrete representation of a process entity
type Entity struct {
	Type string `json:type`
	Code string `json:code`
}

func DecodeProcess(b []byte) (p *Process, err error) {
	err = json.Unmarshal(b, &p)
	if err != nil {
		return
	}
	return
}

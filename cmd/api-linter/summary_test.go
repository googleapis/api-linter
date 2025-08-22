// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/v2/lint"
)

func TestCreateSummary(t *testing.T) {
	tests := []struct {
		name        string
		data        []lint.Response
		wantSummary map[string]map[string]int
	}{{
		name:        "Empty input",
		data:        []lint.Response{},
		wantSummary: make(map[string]map[string]int),
	}, {
		name: "Example with a couple of responses",
		data: []lint.Response{
			{
				FilePath: "example.proto",
				Problems: []lint.Problem{
					{RuleID: "core::naming_formats::field_names"},
					{RuleID: "core::naming_formats::field_names"},
				},
			},
			{
				FilePath: "example2.proto",
				Problems: []lint.Problem{
					{RuleID: "core::0131::request_message::name"},
					{RuleID: "core::0132::response_message::name"},
				},
			},
			{
				FilePath: "example3.proto",
				Problems: []lint.Problem{
					{RuleID: "core::naming_formats::field_names"},
				},
			},
			{
				FilePath: "example4.proto",
				Problems: []lint.Problem{
					{RuleID: "core::naming_formats::field_names"},
					{RuleID: "core::0132::response_message::name"},
				},
			},
		},
		wantSummary: map[string]map[string]int{
			"core::0131::request_message::name": {"example2.proto": 1},
			"core::0132::response_message::name": {
				"example2.proto": 1,
				"example4.proto": 1,
			},
			"core::naming_formats::field_names": {
				"example.proto":  2,
				"example3.proto": 1,
				"example4.proto": 1,
			},
		},
	}}
	for _, test := range tests {
		want := test.wantSummary
		got := createSummary(test.data)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("createSummary() mismatch (-want +got):\n%s", diff)
		}
	}
}

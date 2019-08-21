package main

import (
	"reflect"
	"testing"

	"github.com/googleapis/api-linter/lint"
)


func TestCreateSummary(t *testing.T) {
	tests := []struct {
		name	string
		data	[]lint.Response
		wantSummary	*LintSummary
	}{{
		name: "Empty input",
		data: []lint.Response{},
		wantSummary: &LintSummary{violationData: make(map[string]map[string]bool)},
	},
	{
		name: "Example with a couple of responses",
		data: []lint.Response{{
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
		}},
		wantSummary: &LintSummary{
			violationData: map[string]map[string]bool{
				"core::0131::request_message::name": map[string]bool{"example2.proto": true},
				"core::0132::response_message::name": map[string]bool{
					"example2.proto": true,
					"example4.proto": true,
				},
				"core::naming_formats::field_names": map[string]bool{
					"example.proto": true,
					"example3.proto": true,
					"example4.proto": true,
				},

			},
			longestRuleLen: len("core::0132::response_message::name"),
			numSourceFiles: 4,
		},
		//	map[string]int{
		//	"core::0131::request_message::name": 1,
		//	"core::0132::response_message::name": 2,
		//	"core::naming_formats::field_names": 3,
		//},
	}}
	for _, test := range tests {
		gotSummary := createSummary(test.data)
		if !reflect.DeepEqual(gotSummary.violationData, test.wantSummary.violationData) {
			t.Errorf("Incorrect violation data:\nGot: %#v\n Want: %#v", gotSummary.violationData, test.wantSummary.violationData)
		}
		if gotSummary.longestRuleLen != test.wantSummary.longestRuleLen {
			t.Errorf("Incorrect longest rule length:\nGot: %d want: %d", gotSummary.longestRuleLen, test.wantSummary.longestRuleLen)
		}
		if gotSummary.numSourceFiles != test.wantSummary.numSourceFiles {
			t.Errorf("Incorrect numSourceFiles:\nGot: %d want: %d", gotSummary.numSourceFiles, test.wantSummary.numSourceFiles)
		}
	}
}
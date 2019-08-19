package main

import (
	"reflect"
	"testing"

	"github.com/googleapis/api-linter/lint"
)


func TestSummaryCli(t *testing.T) {
	tests := []struct {
		name	string
		data	[]lint.Response
		wantSummary	map[string]int
		wantLongestRuleLen int
	}{{
		name: "Empty input",
		data: []lint.Response{},
		wantSummary: make(map[string]int),
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
		wantSummary: map[string]int{
			"core::0131::request_message::name": 1,
			"core::0132::response_message::name": 2,
			"core::naming_formats::field_names": 3,
		},
		wantLongestRuleLen: 34,
	}}
	for _, test := range tests {
		gotSummary, gotLongestRuleLen := createSummary(test.data)
		if !reflect.DeepEqual(gotSummary, test.wantSummary) {
			t.Errorf("Incorrect Summary Output: \nGot: %#v\n Want: %#v", gotSummary, test.wantSummary)
		}
		if gotLongestRuleLen != test.wantLongestRuleLen {
			t.Errorf("Incorrect")
		}
	}
}
package main

import (
	"reflect"
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestCreateSummary(t *testing.T) {
	tests := []struct {
		name        string
		data        []lint.Response
		wantSummary *LintSummary
	}{{
		name:        "Empty input",
		data:        []lint.Response{},
		wantSummary: &LintSummary{Violations: make(map[string]map[string]bool)},
	}, {
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
			Violations: map[string]map[string]bool{
				"core::0131::request_message::name": map[string]bool{"example2.proto": true},
				"core::0132::response_message::name": map[string]bool{
					"example2.proto": true,
					"example4.proto": true,
				},
				"core::naming_formats::field_names": map[string]bool{
					"example.proto":  true,
					"example3.proto": true,
					"example4.proto": true,
				},
			},
			LongestRuleLen: len("core::0132::response_message::name"),
			NumSourceFiles: 4,
		},
	}}
	for _, test := range tests {
		gotSummary := createSummary(test.data)
		if !reflect.DeepEqual(gotSummary.Violations, test.wantSummary.Violations) {
			t.Errorf("Incorrect violation data:\nGot: %#v\n Want: %#v", gotSummary.Violations, test.wantSummary.Violations)
		}
		if gotSummary.LongestRuleLen != test.wantSummary.LongestRuleLen {
			t.Errorf("Incorrect longest rule length:\nGot: %d want: %d", gotSummary.LongestRuleLen, test.wantSummary.LongestRuleLen)
		}
		if gotSummary.NumSourceFiles != test.wantSummary.NumSourceFiles {
			t.Errorf("Incorrect numSourceFiles:\nGot: %d want: %d", gotSummary.NumSourceFiles, test.wantSummary.NumSourceFiles)
		}
	}
}

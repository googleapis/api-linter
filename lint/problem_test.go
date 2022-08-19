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

package lint

import (
	"encoding/json"
	"strings"
	"testing"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/builder"
	"gopkg.in/yaml.v2"
)

func TestProblemJSON(t *testing.T) {
	problem := &Problem{
		Message:  "foo bar",
		Location: &dpb.SourceCodeInfo_Location{Span: []int32{2, 0, 42}},
		RuleID:   "core::0131",
	}
	serialized, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Could not marshal Problem to JSON.")
	}
	tests := []struct {
		testName string
		token    string
	}{
		{"Message", `"message":"foo bar"`},
		{"LineNumber", `"line_number":3`},
		{"ColumnNumberStart", `"column_number":1`},
		{"ColumnNumberEnd", `"column_number":42`},
		{"RuleID", `"rule_id":"core::0131"`},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if !strings.Contains(string(serialized), test.token) {
				t.Errorf("Got\n%v\nExpected `%s` to be present.", string(serialized), test.token)
			}
		})
	}
}

func TestProblemYAML(t *testing.T) {
	problem := &Problem{
		Message:  "foo bar",
		Location: &dpb.SourceCodeInfo_Location{Span: []int32{2, 0, 5, 70}},
		RuleID:   "core::0131",
	}
	serialized, err := yaml.Marshal(problem)
	if err != nil {
		t.Fatalf("Could not marshal Problem to YAML.")
	}
	tests := []struct {
		testName string
		token    string
	}{
		{"Message", `message: foo bar`},
		{"LineNumberStart", `line_number: 3`},
		{"LintNumberEnd", `line_number: 6`},
		{"ColumnNumberStart", `column_number: 1`},
		{"ColumnNumberEnd", `column_number: 70`},
		{"RuleID", `rule_id: core::0131`},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if !strings.Contains(string(serialized), test.token) {
				t.Errorf("Got\n%v\nExpected `%s` to be present.", string(serialized), test.token)
			}
		})
	}
}

func TestProblemDescriptor(t *testing.T) {
	m, err := builder.NewMessage("Foo").Build()
	if err != nil {
		t.Fatalf("%v", err)
	}
	m.GetSourceInfo().Span = []int32{42, 0, 79}
	problem := &Problem{
		Message:    "foo bar",
		Descriptor: m,
		RuleID:     "core::0131",
	}
	serialized, err := yaml.Marshal(problem)
	if err != nil {
		t.Fatalf("Could not marshal Problem to YAML.")
	}
	tests := []struct {
		testName string
		token    string
	}{
		{"Message", `message: foo bar`},
		{"LineNumber", `line_number: 43`},
		{"ColumnNumberStart", `column_number: 1`},
		{"ColumnNumberEnd", `column_number: 79`},
		{"RuleID", `rule_id: core::0131`},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if !strings.Contains(string(serialized), test.token) {
				t.Errorf("Got\n%v\nExpected `%s` to be present.", string(serialized), test.token)
			}
		})
	}
}

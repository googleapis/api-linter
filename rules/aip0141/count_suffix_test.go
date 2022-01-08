// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0141

import (
	"testing"

	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/testutils"
)

func TestCount(t *testing.T) {
	tests := []struct {
		testName  string
		FieldName string
		want      string
		err       bool
	}{
		{"Valid", "item_count", "item_count", false},
		{"NumItems", "num_items", "item_count", true},
		{"NumResults", "num_results", "result_count", true},
		{"NumErrors", "num_errors", "error_count", true},
		{"NumMoose", "num_moose", "moose_count", true}, // singular of "moose" is "moose"
		{"NumGeese", "num_geese", "goose_count", true}, // singular of "geese" is "goose"
	}
	for _, test := range tests {
		file := testutils.ParseProto3Tmpl(t, `
			message Job {
				int32 {{.FieldName}} = 1;
			}
		`, test)
		field := file.GetMessageTypes()[0].GetFields()[0]
		wantProblems := testutils.Problems{}
		if test.err {
			wantProblems = append(wantProblems, lint.Problem{
				Message:    "_count suffix",
				Suggestion: test.want,
				Descriptor: field,
			})
		}
		if diff := wantProblems.Diff(count.Lint(file)); diff != "" {
			t.Errorf(diff)
		}
	}
}

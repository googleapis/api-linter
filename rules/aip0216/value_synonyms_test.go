// Copyright 2020 Google LLC
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

package aip0216

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestValueSynonyms(t *testing.T) {
	for _, test := range []struct {
		name      string
		EnumName  string
		ValueName string
		problems  testutils.Problems
	}{
		{"ValidSucceeded", "State", "SUCCEEDED", nil},
		{"ValidDeleted", "State", "DELETED", nil},
		{"InvalidSuccess", "State", "SUCCESS", testutils.Problems{{Suggestion: "SUCCEEDED"}}},
		{"InvalidSuccessful", "State", "SUCCESSFUL", testutils.Problems{{Suggestion: "SUCCEEDED"}}},
		{"InvalidCanceled", "State", "CANCELED", testutils.Problems{{Suggestion: "CANCELLED"}}},
		{"InvalidCanceling", "State", "CANCELING", testutils.Problems{{Suggestion: "CANCELLING"}}},
		{"InvalidFail", "State", "FAIL", testutils.Problems{{Suggestion: "FAILED"}}},
		{"InvalidFailure", "State", "FAILURE", testutils.Problems{{Suggestion: "FAILED"}}},
		{"Irrelevant", "Foo", "SUCCESSFUL", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				enum {{.EnumName}} {
					{{.EnumName}}_UNSPECIFIED = 0;
					{{.ValueName}} = 1;
				}
			`, test)
			ev := file.GetEnumTypes()[0].GetValues()[1]
			if diff := test.problems.SetDescriptor(ev).Diff(valueSynonyms.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

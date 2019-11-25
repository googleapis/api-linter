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

package aip0217

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/stoewer/go-strcase"
)

func TestSynonyms(t *testing.T) {
	for _, test := range []struct {
		FieldName string
		problems  testutils.Problems
	}{
		{"unreachable", testutils.Problems{}},
		{"unreachable_locations", testutils.Problems{{Suggestion: "unreachable"}}},
	} {
		t.Run(strcase.UpperCamelCase(test.FieldName), func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message ListBooksResponse {
					repeated Book books = 1;
					string next_page_token = 2;
					repeated string {{.FieldName}} = 3;
				}
				message Book {}

				// Ignored because next_page_token is not present.
				message Irrelevant {
					repeated string unreachable_locations = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[2]
			if diff := test.problems.SetDescriptor(field).Diff(synonyms.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

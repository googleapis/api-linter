// Copyright 2020 Google LLC
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

package aip0140

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestURI(t *testing.T) {
	for _, test := range []struct {
		name     string
		Comment  string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "", "uri", nil},
		{"ValidPrefix", "", "uri_foo", nil},
		{"ValidSuffix", "", "foo_uri", nil},
		{"ValidIntermixed", "", "foo_uri_bar", nil},
		{"ValidURL", "", "url", nil},
		{"ValidURLPrefix", "", "url_foo", nil},
		{"ValidURLSuffix", "", "foo_url", nil},
		{"ValidURLIntermixed", "", "foo_url_bar", nil},
		{"Invalid", "// A URI.", "url", testutils.Problems{{Message: "uri", Suggestion: "uri"}}},
		{"InvalidPrefix", "// A uri.", "url_foo", testutils.Problems{{Message: "uri", Suggestion: "uri_foo"}}},
		{"InvalidSuffix", "// A URI.", "foo_url", testutils.Problems{{Message: "uri", Suggestion: "foo_uri"}}},
		{"InvalidIntermixed", "// A uri.", "foo_url_bar", testutils.Problems{{Message: "uri", Suggestion: "foo_uri_bar"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Foo {
				  {{.Comment}}
					string {{.Field}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(uri.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

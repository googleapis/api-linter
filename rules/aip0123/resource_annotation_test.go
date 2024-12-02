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

package aip0123

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestResourceAnnotation(t *testing.T) {
	// The rule should pass if the option is present on a resource message.
	t.Run("Present", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
				string name = 1;
			}
		`)
		if diff := (testutils.Problems{}).Diff(resourceAnnotation.Lint(f)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("SkipNested", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			message Foo {
				message Bar {
					string name = 1;
				}
				Bar bar = 1;
			}
		`)
		if diff := (testutils.Problems{}).Diff(resourceAnnotation.Lint(f)); diff != "" {
			t.Error(diff)
		}
	})

	// The rule should fail if the option is absent on a resource message,
	// but pass on messages that are not resource messages.
	for _, test := range []struct {
		name        string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"ValidNoNameField", "Book", "title", testutils.Problems{}},
		{"ValidRequestMessage", "GetBookRequest", "name", testutils.Problems{}},
		{"ValidResponseMessage", "GetBookResponse", "name", testutils.Problems{}},
		{"Invalid", "Book", "name", testutils.Problems{{Message: "google.api.resource"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					string {{.FieldName}} = 1;
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(resourceAnnotation.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

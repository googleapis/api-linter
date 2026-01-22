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

package aip0143

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestFieldNames(t *testing.T) {
	tests := []struct {
		FieldName         string
		ResourceReference string
		problems          testutils.Problems
	}{
		{
			FieldName: "something_random",
			problems:  testutils.Problems{},
		},
		{
			FieldName: "content_type",
			problems:  testutils.Problems{{Suggestion: "mime_type"}},
		},
		{
			FieldName: "country",
			problems:  testutils.Problems{{Suggestion: "region_code"}},
		},
		{
			FieldName: "country_code",
			problems:  testutils.Problems{{Suggestion: "region_code"}},
		},
		{
			FieldName: "region_code",
			problems:  testutils.Problems{},
		},
		{
			FieldName: "currency",
			problems:  testutils.Problems{{Suggestion: "currency_code"}},
		},
		{
			FieldName: "currency_code",
			problems:  testutils.Problems{},
		},
		{
			FieldName: "language",
			problems:  testutils.Problems{{Suggestion: "language_code"}},
		},
		{
			FieldName: "language_code",
			problems:  testutils.Problems{},
		},
		{
			FieldName: "mime",
			problems:  testutils.Problems{{Suggestion: "mime_type"}},
		},
		{
			FieldName: "mimetype",
			problems:  testutils.Problems{{Suggestion: "mime_type"}},
		},
		{
			FieldName: "mime_type",
			problems:  testutils.Problems{},
		},
		{
			FieldName: "timezone",
			problems:  testutils.Problems{{Suggestion: "time_zone"}},
		},
		{
			FieldName: "time_zone",
			problems:  testutils.Problems{},
		},
		{
			// Skip when field represents a resource name.
			FieldName:         "language",
			ResourceReference: `[(google.api.resource_reference).type = "example.com/Language"]`,
			problems:          testutils.Problems{},
		},
	}
	for _, test := range tests {
		t.Run(test.FieldName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Foo {
					string {{.FieldName}} = 1 {{.ResourceReference}};
				}
			`, test)
			field := file.Messages().Get(0).Fields().Get(0)
			problems := fieldNames.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

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

package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

var (
	generatedURI                       = "string generated_uri = 1;"
	generatedURIWithOutputOnlyBehavior = `string generated_uri = 1 [
	(google.api.field_behavior) = OUTPUT_ONLY];`
)

func TestOutput(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "Valid",
			comment:  "Output Only",
			field:    generatedURIWithOutputOnlyBehavior,
			problems: nil,
		},
		{
			name:     "Valid",
			comment:  "@OutputOnly",
			field:    generatedURIWithOutputOnlyBehavior,
			problems: nil,
		},
		{
			name:     "Valid-exclude-etag",
			comment:  "Output Only.",
			field:    "string etag = 1;",
			problems: nil,
		},
		{
			name:    "Invalid-Output Only",
			comment: "Output Only",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-OutputOnly",
			comment: "OutputOnly",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@OutputOnly",
			comment: "@OutputOnly",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-output_only",
			comment: "output_only",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@output_only",
			comment: "@output_only",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-OUTPUT_ONLY",
			comment: "OUTPUT_ONLY",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-outputOnly",
			comment: "outputOnly",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@outputOnly",
			comment: "@outputonly",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@OUTPUT_ONLY",
			comment: "@OUTPUT_ONLY",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-output_only_free_text",
			comment: "This field is output only",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-!outputOnly",
			comment: "!outputOnly",
			field:   generatedURI,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
				import "google/api/field_behavior.proto";
				message Book {
					// Secret to be stored in the book.
					// {{.Comment}}
					{{.Field}}
				}`
			file := testutils.ParseProto3Tmpl(t, template,
				struct {
					Comment string
					Field   string
				}{test.comment, test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := outputOnly.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

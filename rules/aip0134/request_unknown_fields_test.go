// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package aip0134

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		MessageName string
		FieldName   string
		FieldType   string
		problems    testutils.Problems
	}{
		// Use BigBook instead of Book to test correct casing logic
		{
			"UpdateMask", "UpdateBigBookRequest", "update_mask",
			"google.protobuf.FieldMask",
			testutils.Problems{},
		},
		{
			"ValidateOnly", "UpdateBigBookRequest", "validate_only",
			"bool",
			testutils.Problems{},
		},
		{
			"View", "UpdateBigBookRequest", "view",
			"BigBookView",
			testutils.Problems{},
		},
		{
			"SuffixedView", "UpdateBigBookRequest", "custom_view",
			"BigBookView",
			testutils.Problems{},
		},
		{
			"Invalid", "UpdateBigBookRequest", "application_id",
			"string",
			testutils.Problems{{Message: "Unexpected field"}},
		},
		{
			"InvalidCasing", "UpdateBigBookRequest", "bigbook",
			"string",
			testutils.Problems{{Message: "Unexpected field"}},
		},
		{
			"Irrelevant", "AcquireBigBookRequest", "application_id",
			"string",
			testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/protobuf/field_mask.proto";
				enum BigBookView {
					BIG_BOOK_VIEW_UNSPECIFIED = 0;
					BIG_BOOK_VIEW_BASIC = 1;
					BIG_BOOK_VIEW_FULL = 2;
				}
				message {{.MessageName}} {
					BigBook big_book = 1;
					{{.FieldType}} {{.FieldName}} = 2;
				}
				message BigBook {}
			`, test)

			// Run the lint rule, and establish that it returns the correct problems.
			problems := unknownFields.Lint(f)
			if diff := test.problems.SetDescriptor(f.Messages().Get(0).Fields().Get(1)).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

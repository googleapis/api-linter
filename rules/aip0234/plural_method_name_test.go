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

package aip0234

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestMethodPluralResourceName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		MethodName    string
		URISuffix     string
		ResponseItems string
		problems      testutils.Problems
	}{
		{
			testName:      "Valid-BatchUpdateBooks",
			MethodName:    "BatchUpdateBooks",
			URISuffix:     "books:batchUpdate",
			ResponseItems: "repeated Book books = 1;",
			problems:      testutils.Problems{},
		},
		{
			testName:      "Valid-BatchUpdateMen",
			MethodName:    "BatchUpdateMen",
			URISuffix:     "men:batchUpdate",
			ResponseItems: "repeated Other men = 1;",
			problems:      testutils.Problems{},
		},

		{
			testName:      "ValidBatchUpdateTitles-NonMessageItems",
			MethodName:    "BatchUpdateTitles",
			URISuffix:     "titles:batchUpdate",
			ResponseItems: "repeated string titles = 1;",
			problems:      testutils.Problems{},
		},
		{
			testName:      "Invalid-SingularBus",
			MethodName:    "BatchUpdateBus",
			URISuffix:     "bus:batchUpdate",
			ResponseItems: "repeated Other buses = 1;",
			problems: testutils.Problems{{
				Suggestion: "BatchUpdateBuses",
			}},
		},
		{
			testName:      "Invalid-SingularCorpPerson",
			MethodName:    "BatchUpdateCorpPerson",
			URISuffix:     "corpPerson:batchUpdate",
			ResponseItems: "repeated Other corp_people = 1;",
			problems: testutils.Problems{{
				Suggestion: "BatchUpdateCorpPeople",
			}},
		},
		{
			testName:   "Skip-NotBatchUpdate",
			MethodName: "AcquireBook",
			URISuffix:  "book",
			problems:   testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";

				service BookService {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							post: "/v1/{parent=publishers/*}/{{.URISuffix}}"
							body: "*"
						};
					}
				}

				message {{.MethodName}}Request {}

				message {{.MethodName}}Response {
					{{ .ResponseItems }}
				}

				message Book {
				  option (google.api.resource) = {
				    type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
					singular: "book"
					plural: "books"
				  };
				}

				message Other {}
				`, test)

			m := file.Services().Get(0).Methods().Get(0)

			problems := pluralMethodName.Lint(file)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

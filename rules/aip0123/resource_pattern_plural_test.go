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

func TestResourcePatternPluralSimple(t *testing.T) {
	for _, test := range []struct {
		name     string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "publishers/{publisher}/bookShelves/{book_shelf}", testutils.Problems{}},
		{"Invalid", "publishers/{publisher}/bookShelfs/{book_shelf}", testutils.Problems{{Message: "collection segment must be the resource plural"}}},
		{"SkipRootLevel", "publishers/{publisher}", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				message BookShelf {
					option (google.api.resource) = {
						type: "library.googleapis.com/BookShelf"
						singular: "bookShelf"
						plural: "bookShelves"
						pattern: "{{.Pattern}}"
					};
					string name = 1;
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(resourcePatternPlural.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestResourcePatternPluralNested(t *testing.T) {
	for _, test := range []struct {
		name          string
		FirstPattern  string
		SecondPattern string
		problems      testutils.Problems
	}{
		{
			name:          "Valid",
			FirstPattern:  "publishers/{publisher}/credits/{credit}",
			SecondPattern: "authors/{author}/credits/{credit}",
			problems:      testutils.Problems{},
		},
		{
			name:          "ValidFull",
			FirstPattern:  "publishers/{publisher}/publisherCredits/{publisher_credit}",
			SecondPattern: "authors/{author}/publisherCredits/{publisher_credit}",
			problems:      testutils.Problems{},
		},
		{
			name:          "InvalidSecondWithFirstNestedName",
			FirstPattern:  "publishers/{publisher}/credits/{credit}",
			SecondPattern: "authors/{author}/publisherCredits/{credit}",
			problems:      testutils.Problems{{Message: `collection segment must be the resource plural "/credits/"`}},
		},
		{
			name:          "InvalidFirstWithReducedSecond",
			FirstPattern:  "publishers/{publisher}/pubCredits/{credit}",
			SecondPattern: "authors/{author}/credits/{credit}",
			problems:      testutils.Problems{{Message: `collection segment must be the resource plural "/credits/"`}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";

				message PublisherCredit {
					option (google.api.resource) = {
						type: "library.googleapis.com/PublisherCredit"
						singular: "publisherCredit"
						plural: "publisherCredits"
						pattern: "{{.FirstPattern}}"
						pattern: "{{.SecondPattern}}"
					};
					string name = 1;
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(resourcePatternPlural.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

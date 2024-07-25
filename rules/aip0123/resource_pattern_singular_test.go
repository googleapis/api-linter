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

func TestResourcePatternSingularSimple(t *testing.T) {
	for _, test := range []struct {
		name     string
		Pattern  string
		problems testutils.Problems
	}{
		{"Valid", "publishers/{publisher}/bookShelves/{book_shelf}", testutils.Problems{}},
		{"ValidSingleton", "publishers/{publisher}/bookShelf", testutils.Problems{}},
		{"Invalid", "publishers/{publisher}/bookShelves/{shelf}", testutils.Problems{{Message: "final segment must include the resource singular"}}},
		{"InvalidSingleton", "publishers/{publisher}/shelf", testutils.Problems{{Message: "final segment must include the resource singular"}}},
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
			if diff := test.problems.SetDescriptor(m).Diff(resourcePatternSingular.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestResourcePatternSingularNested(t *testing.T) {
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
			name:          "ValidSingleton",
			FirstPattern:  "publishers/{publisher}/credit",
			SecondPattern: "authors/{author}/credit",
			problems:      testutils.Problems{},
		},
		{
			name:          "ValidSingletonFull",
			FirstPattern:  "publishers/{publisher}/publisherCredit",
			SecondPattern: "authors/{author}/publisherCredit",
			problems:      testutils.Problems{},
		},
		{
			name:          "InvalidSecondWithFirstNestedName",
			FirstPattern:  "publishers/{publisher}/credits/{credit}",
			SecondPattern: "authors/{author}/credits/{published}",
			problems:      testutils.Problems{{Message: `final segment must include the resource singular "{credit}"`}},
		},
		{
			name:          "InvalidFirstWithReducedSecond",
			FirstPattern:  "publishers/{publisher}/credits/{published}",
			SecondPattern: "authors/{author}/credits/{credit}",
			problems:      testutils.Problems{{Message: `final segment must include the resource singular "{credit}"`}},
		},
		{
			name:          "InvalidSingletonFirstPattern",
			FirstPattern:  "publishers/{publisher}/published",
			SecondPattern: "authors/{author}/credit",
			problems:      testutils.Problems{{Message: `final segment must include the resource singular "credit"`}},
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
			if diff := test.problems.SetDescriptor(m).Diff(resourcePatternSingular.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

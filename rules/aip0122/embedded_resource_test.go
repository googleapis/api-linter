// Copyright 2022 Google LLC
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

package aip0122

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestEmbeddedResource(t *testing.T) {
	for _, test := range []struct {
		name           string
		FieldA, FieldB string
		problems       testutils.Problems
	}{
		{"Valid", "", "", nil},
		{
			"ValidReferences",
			`string library = 2 [(google.api.resource_reference).type = "library.googleapis.com/Library"];`,
			`string librarian = 3 [(google.api.resource_reference).type = "library.googleapis.com/Librarian"];`,
			nil,
		},
		{
			"InvalidEmbeddedResources",
			"Library library = 2;",
			"Librarian librarian = 3;",
			testutils.Problems{
				{
					Message:    "not by embedding",
					Suggestion: `string library = 2 [(google.api.resource_reference).type = "library.googleapis.com/Library"];`,
				},
				{
					Message:    "not by embedding",
					Suggestion: `string librarian = 3 [(google.api.resource_reference).type = "library.googleapis.com/Librarian"];`,
				},
			},
		},
		{
			"InvalidRepeatedEmbeddedResource",
			"repeated Library library = 2;",
			"",
			testutils.Problems{{
				Message:    "not by embedding",
				Suggestion: `repeated string library = 2 [(google.api.resource_reference).type = "library.googleapis.com/Library"];`,
			}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
				string name = 1;

				{{.FieldA}}

				{{.FieldB}}
			}

			message Library {
				option (google.api.resource) = {
					type: "library.googleapis.com/Library"
					pattern: "libraries/{library}"
				};
				string name = 1;
			}

			message Librarian {
				option (google.api.resource) = {
					type: "library.googleapis.com/Librarian"
					pattern: "libraries/{library}/librarians/{librarian}"
				};
				string name = 1;
			}
		`, test)
			m := f.Messages().ByName("Book")

			want := test.problems
			if len(want) > 0 {
				want[0].Descriptor = m.Fields().ByName("library")
			}
			if len(want) > 1 {
				want[1].Descriptor = m.Fields().ByName("librarian")
			}
			if diff := want.Diff(embeddedResource.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEmbeddedResource_Revisions(t *testing.T) {
	for _, test := range []struct {
		name         string
		SnapshotType string
		problems     testutils.Problems
	}{
		{"Valid", "Book", nil},
		{
			"InvalidEmbeddedResource",
			"Library",
			testutils.Problems{{
				Message:    "not by embedding",
				Suggestion: `string snapshot = 2 [(google.api.resource_reference).type = "library.googleapis.com/Library"];`,
			}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
				string name = 1;
			}

			message BookRevision {
				option (google.api.resource) = {
					type: "library.googleapis.com/BookRevision"
					pattern: "publishers/{publisher}/books/{book}/revisions/{revision}"
				};
				string name = 1;

				{{.SnapshotType}} snapshot = 2;
			}

			message Library {
				option (google.api.resource) = {
					type: "library.googleapis.com/Library"
					pattern: "libraries/{library}"
				};
				string name = 1;
			}
		`, test)
			field := f.Messages().ByName("BookRevision").Fields().ByName("snapshot")
			if diff := test.problems.SetDescriptor(field).Diff(embeddedResource.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

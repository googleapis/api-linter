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

package utils

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		pluralizedWord string
	}{
		{"PluralizeSingularWord", "person", "people"},
		{"PluralizePluralWord", "people", "people"},
		{"PluralizeNonstandardPluralWord", "persons", "people"},
		{"PluralizeNoPluralFormWord", "moose", "moose"},
		{"PluralizePluralLatinWord", "cacti", "cacti"},
		{"PluralizeNonstandardPluralLatinWord", "cactuses", "cacti"},
		{"PluralizePluralCamelCaseWord", "student_profiles", "student_profiles"},
		{"PluralizeSingularCamelCaseWord", "student_profile", "student_profiles"},
	}
	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			if got := ToPlural(test.word); got != test.pluralizedWord {
				t.Errorf("Plural(%s) got %s, but want %s", test.word, got, test.pluralizedWord)
			}
		})
	}
}

func TestResourceSingular(t *testing.T) {
	tests := []struct {
		testName   string
		pluralName string
		src        string
		want       string
	}{
		{
			testName:   "AnnotationInSameFile",
			pluralName: "ImpressionMetadata",
			src: `
				import "google/api/resource.proto";

				message BatchUpdateImpressionMetadataRequest {
					string parent = 1;
				}
				message ImpressionMetadata {
					option (google.api.resource) = {
						type: "example.com/ImpressionMetadata"
						pattern: "dataProviders/{dp}/impressionMetadata/{im}"
						singular: "impressionMetadata"
						plural: "impressionMetadata"
					};
				}
			`,
			want: "ImpressionMetadata",
		},
		{
			testName:   "FallbackToGoPluralizeBooks",
			pluralName: "Books",
			src: `
				message BatchUpdateBooksRequest {
					string parent = 1;
				}
			`,
			want: "Book",
		},
		{
			testName:   "UncountableNounPluralEqualsSingular",
			pluralName: "Metadata",
			src: `
				import "google/api/resource.proto";

				message BatchUpdateMetadataRequest {
					string parent = 1;
				}
				message Metadata {
					option (google.api.resource) = {
						type: "example.com/Metadata"
						pattern: "items/{item}/metadata/{metadata}"
						singular: "metadata"
						plural: "metadata"
					};
				}
			`,
			want: "Metadata",
		},
		{
			testName:   "MessageNameMatchesPluralName",
			pluralName: "CursorData",
			src: `
				import "google/api/resource.proto";

				message BatchUpdateCursorDataRequest {
					string parent = 1;
				}
				message CursorData {
					option (google.api.resource) = {
						type: "example.com/CursorData"
						pattern: "items/{item}/cursorData/{cursor_data}"
						singular: "cursorDatum"
						plural: "cursorData"
					};
				}
			`,
			want: "CursorDatum",
		},
		{
			testName:   "NoAnnotationLatinWord",
			pluralName: "Data",
			src: `
				message BatchUpdateDataRequest {
					string parent = 1;
				}
			`,
			want: "Datum",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3String(t, test.src)
			m := file.Messages().Get(0)
			got := ResourceSingular(test.pluralName, m)
			if got != test.want {
				t.Errorf("ResourceSingular(%q) = %q, want %q", test.pluralName, got, test.want)
			}
		})
	}
}

func TestResourceSingularImportedFile(t *testing.T) {
	// Verify that ResourceSingular finds the resource annotation in a
	// directly imported file, not just the same file.
	files := testutils.ParseProtoStrings(t, map[string]string{
		"resource.proto": `
			syntax = "proto3";
			import "google/api/resource.proto";

			message ImpressionMetadata {
				option (google.api.resource) = {
					type: "example.com/ImpressionMetadata"
					pattern: "dataProviders/{dp}/impressionMetadata/{im}"
					singular: "impressionMetadata"
					plural: "impressionMetadata"
				};
			}
		`,
		"service.proto": `
			syntax = "proto3";
			import "resource.proto";

			message BatchUpdateImpressionMetadataRequest {
				string parent = 1;
			}
		`,
	})

	serviceFile := files["service.proto"]
	m := serviceFile.Messages().Get(0)
	got := ResourceSingular("ImpressionMetadata", m)
	if got != "ImpressionMetadata" {
		t.Errorf("ResourceSingular(\"ImpressionMetadata\") = %q, want \"ImpressionMetadata\"", got)
	}
}

func TestResourceSingularTransitiveImport(t *testing.T) {
	// Verify that ResourceSingular finds the resource annotation
	// through transitive imports (service -> common -> resource).
	files := testutils.ParseProtoStrings(t, map[string]string{
		"resource.proto": `
			syntax = "proto3";
			import "google/api/resource.proto";

			message ImpressionMetadata {
				option (google.api.resource) = {
					type: "example.com/ImpressionMetadata"
					pattern: "dataProviders/{dp}/impressionMetadata/{im}"
					singular: "impressionMetadata"
					plural: "impressionMetadata"
				};
			}
		`,
		"common.proto": `
			syntax = "proto3";
			import "resource.proto";

			message CommonFields {
				string parent = 1;
			}
		`,
		"service.proto": `
			syntax = "proto3";
			import "common.proto";

			message BatchUpdateImpressionMetadataRequest {
				string parent = 1;
			}
		`,
	})

	serviceFile := files["service.proto"]
	m := serviceFile.Messages().Get(0)
	got := ResourceSingular("ImpressionMetadata", m)
	if got != "ImpressionMetadata" {
		t.Errorf("ResourceSingular(\"ImpressionMetadata\") = %q, want \"ImpressionMetadata\"", got)
	}
}

// Copyright 2020 Google LLC
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

package aip0124

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestValidReference(t *testing.T) {
	for _, test := range []struct {
		name     string
		Field    string
		Ref      string
		problems testutils.Problems
	}{
		{"Found", "type", "library.googleapis.com/Book", nil},
		{"FoundChild", "child_type", "library.googleapis.com/Book", nil},
		{"NotFound", "type", "library.googleapis.com/Bok", testutils.Problems{{Message: "Could not find"}}},
		{"NotFoundChild", "child_type", "library.googleapis.com/Bok", testutils.Problems{{Message: "Could not find"}}},
		{"Ignored", "type", "cloudresourcemanager.googleapis.com/Project", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("SameFileResourceDefinition", func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					import "google/api/resource.proto";
					option (google.api.resource_definition) = {
						type: "library.googleapis.com/Book"
					};
					message Foo {
						string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
						string irrelevant = 2;
					}
				`, test)
				field := f.GetMessageTypes()[0].GetFields()[0]
				if diff := test.problems.SetDescriptor(field).Diff(validReference.Lint(f)); diff != "" {
					t.Errorf(diff)
				}
			})

			t.Run("SameFileResourceMessage", func(t *testing.T) {
				f := testutils.ParseProto3Tmpl(t, `
					import "google/api/resource.proto";
					message Book {
						option (google.api.resource) = {
							type: "library.googleapis.com/Book"
						};
					}
					message Foo {
						string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
					}
				`, test)
				field := f.GetMessageTypes()[1].GetFields()[0]
				if diff := test.problems.SetDescriptor(field).Diff(validReference.Lint(f)); diff != "" {
					t.Errorf(diff)
				}
			})
		})

		t.Run("DirectDependency", func(t *testing.T) {
			files := testutils.ParseProto3Tmpls(t, map[string]string{
				"dep.proto": `
					import "google/api/resource.proto";
					message Book {
						option (google.api.resource) = {
							type: "library.googleapis.com/Book"
						};
					}
				`,
				"leaf.proto": `
					import "google/api/resource.proto";
					import "dep.proto";
					message Foo {
						string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
					}
				`,
			}, test)
			file := files["leaf.proto"]
			field := file.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(validReference.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})

		t.Run("RemoteDependency", func(t *testing.T) {
			files := testutils.ParseProto3Tmpls(t, map[string]string{
				"dep.proto": `
					import "google/api/resource.proto";
					message Book {
						option (google.api.resource) = {
							type: "library.googleapis.com/Book"
						};
					}
				`,
				"intermediate.proto": `import "dep.proto";`,
				"leaf.proto": `
					import "google/api/resource.proto";
					import "intermediate.proto";
					message Foo {
						string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
					}
				`,
			}, test)
			file := files["leaf.proto"]
			field := file.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(validReference.Lint(file)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

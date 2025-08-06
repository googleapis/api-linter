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

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
)

func TestReferenceSamePackage(t *testing.T) {
	for _, test := range []struct {
		name         string
		Field        string
		Ref          string
		OtherPackage string
		problems     testutils.Problems
	}{
		{"SamePkg", "type", "library.googleapis.com/Book", "same", nil},
		{"SamePkgChild", "child_type", "library.googleapis.com/Book", "same", nil},
		{"NotSamePkg", "type", "library.googleapis.com/Book", "other", testutils.Problems{{Message: "same package"}}},
		{"NotSamePkgChild", "child_type", "library.googleapis.com/Book", "other", testutils.Problems{{Message: "same package"}}},
		{"Ignored", "type", "cloudresourcemanager.googleapis.com/Project", "other", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("DirectDependency", func(t *testing.T) {
				files := testutils.ParseProto3Tmpls(t, map[string]string{
					"dep.proto": `
						package {{.OtherPackage}};
						import "google/api/resource.proto";
						message Book {
							option (google.api.resource) = {
								type: "{{.Ref}}"
							};
						}
					`,
					"leaf.proto": `
						package same;
						import "google/api/resource.proto";
						import "dep.proto";
						message Foo {
							string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
						}
					`,
				}, test)
				file := files["leaf.proto"]
				field := file.Messages().Get(0).Fields().Get(0)
				if diff := test.problems.SetDescriptor(field).Diff(referenceSamePackage.Lint(file)); diff != "" {
					t.Error(diff)
				}
			})

			t.Run("DirectDependencyResourceDefinition", func(t *testing.T) {
				files := testutils.ParseProto3Tmpls(t, map[string]string{
					"dep.proto": `
						package {{.OtherPackage}};
						import "google/api/resource.proto";
						option (google.api.resource_definition) = {
							type: "{{.Ref}}"
						};
					`,
					"leaf.proto": `
						package same;
						import "google/api/resource.proto";
						import "dep.proto";
						message Foo {
							string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
						}
					`,
				}, test)
				file := files["leaf.proto"]
				field := file.Messages().Get(0).Fields().Get(0)
				if diff := test.problems.SetDescriptor(field).Diff(referenceSamePackage.Lint(file)); diff != "" {
					t.Error(diff)
				}
			})

			t.Run("RemoteDependency", func(t *testing.T) {
				files := testutils.ParseProto3Tmpls(t, map[string]string{
					"dep.proto": `
						package {{.OtherPackage}};
						import "google/api/resource.proto";
						message Book {
							option (google.api.resource) = {
								type: "library.googleapis.com/Book"
							};
						}
					`,
					"intermediate.proto": `import "dep.proto";`,
					"leaf.proto": `
						package same;
						import "google/api/resource.proto";
						import "intermediate.proto";
						message Foo {
							string book = 1 [(google.api.resource_reference).{{.Field}} = "{{.Ref}}"];
						}
					`,
				}, test)
				file := files["leaf.proto"]
				field := file.Messages().Get(0).Fields().Get(0)
				if diff := test.problems.SetDescriptor(field).Diff(referenceSamePackage.Lint(file)); diff != "" {
					t.Error(diff)
				}
			})
		})
	}
}

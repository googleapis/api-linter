// Copyright 2021 Google LLC
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

func TestDuplicateResource(t *testing.T) {
	f := testutils.ParseProto3Tmpls(t, map[string]string{
		"dep.proto": `
			import "google/api/resource.proto";
			package xyz;
			message Publisher {
				option (google.api.resource) = { type: "library.googleapis.com/Publisher" };
			}
			`,
		"test.proto": `
			import "dep.proto";
			import "google/api/resource.proto";
			package abc;
			option (google.api.resource_definition) = { type: "library.googleapis.com/Publisher" };
			option (google.api.resource_definition) = { type: "library.googleapis.com/Author" };
			option (google.api.resource_definition) = { type: "library.googleapis.com/Editor" };
			message Book {
				option (google.api.resource) = { type: "library.googleapis.com/Book" };
			}
			message Author {
				option (google.api.resource) = { type: "library.googleapis.com/Author" };
			}
			message Foo {
				message Tome {
					option (google.api.resource) = { type: "library.googleapis.com/Book" };
				}
			}`,
	}, nil)["test.proto"]
	want := testutils.Problems{
		{
			Message:    "resource \"library.googleapis.com/Author\": `google.api.resource_definition` 1 in file `test.proto`, message `abc.Author`.",
			Descriptor: f,
		},
		{
			Message:    "resource \"library.googleapis.com/Author\": `google.api.resource_definition` 1 in file `test.proto`, message `abc.Author`.",
			Descriptor: f.Messages().Get(1),
		},
		{
			Message:    "resource \"library.googleapis.com/Book\": message `abc.Book`, message `abc.Foo.Tome`.",
			Descriptor: f.Messages().Get(0),
		},
		{
			Message:    "resource \"library.googleapis.com/Book\": message `abc.Book`, message `abc.Foo.Tome`.",
			Descriptor: f.Messages().Get(2).Messages().Get(0),
		},
		{
			Message:    "resource \"library.googleapis.com/Publisher\": message `xyz.Publisher`.",
			Descriptor: f,
		},
	}
	if diff := want.Diff(duplicateResource.Lint(f)); diff != "" {
		t.Fatal(diff)
	}
}

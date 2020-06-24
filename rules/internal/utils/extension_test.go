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

	"bitbucket.org/creachadair/stringset"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestGetFieldBehavior(t *testing.T) {
	fd := testutils.ParseProto3String(t, `
		import "google/api/field_behavior.proto";

		message Book {
			string name = 1 [
				(google.api.field_behavior) = IMMUTABLE,
				(google.api.field_behavior) = OUTPUT_ONLY];

			string title = 2 [(google.api.field_behavior) = REQUIRED];

			string summary = 3;
		}
	`)
	msg := fd.GetMessageTypes()[0]
	tests := []struct {
		fieldName      string
		fieldBehaviors stringset.Set
	}{
		{"name", stringset.New("IMMUTABLE", "OUTPUT_ONLY")},
		{"title", stringset.New("REQUIRED")},
		{"summary", stringset.New()},
	}
	for _, test := range tests {
		t.Run(test.fieldName, func(t *testing.T) {
			f := msg.FindFieldByName(test.fieldName)
			if diff := cmp.Diff(GetFieldBehavior(f), test.fieldBehaviors); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestGetMethodSignatures(t *testing.T) {
	for _, test := range []struct {
		name       string
		want       [][]string
		Signatures string
	}{
		{"Zero", [][]string{}, ""},
		{"One", [][]string{{"name"}}, `option (google.api.method_signature) = "name";`},
		{"Two",
			[][]string{{"name"}, {"name", "read_mask"}},
			`option (google.api.method_signature) = "name";
			 option (google.api.method_signature) = "name,read_mask";`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/client.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						{{.Signatures}}
					}
				}
				message Book {}
				message GetBookRequest {}
			`, test)
			method := f.GetServices()[0].GetMethods()[0]
			if diff := cmp.Diff(GetMethodSignatures(method), test.want); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestGetOperationInfo(t *testing.T) {
	fd := testutils.ParseProto3String(t, `
		import "google/longrunning/operations.proto";
		service Library {
			rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "WriteBookResponse"
					metadata_type: "WriteBookMetadata"
				};
			}
		}
		message WriteBookRequest {}
	`)
	lro := GetOperationInfo(fd.GetServices()[0].GetMethods()[0])
	if got, want := lro.ResponseType, "WriteBookResponse"; got != want {
		t.Errorf("Response type - got %q, want %q.", got, want)
	}
	if got, want := lro.MetadataType, "WriteBookMetadata"; got != want {
		t.Errorf("Metadata type - got %q, want %q.", got, want)
	}
}

func TestGetOperationInfoNone(t *testing.T) {
	fd := testutils.ParseProto3String(t, `
		service Library {
			rpc GetBook(GetBookRequest) returns (Book);
		}
		message GetBookRequest {}
		message Book {}
	`)
	lro := GetOperationInfo(fd.GetServices()[0].GetMethods()[0])
	if lro != nil {
		t.Errorf("Got %v, expected nil LRO annotation.", lro)
	}
}

func TestGetResource(t *testing.T) {
	t.Run("Present", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
			}
		`)
		resource := GetResource(f.GetMessageTypes()[0])
		if got, want := resource.GetType(), "library.googleapis.com/Book"; got != want {
			t.Errorf("Got %q, expected %q.", got, want)
		}
		if got, want := resource.GetPattern()[0], "publishers/{publisher}/books/{book}"; got != want {
			t.Errorf("Got %q, expected %q.", got, want)
		}
	})
	t.Run("Absent", func(t *testing.T) {
		f := testutils.ParseProto3String(t, "message Book {}")
		if got := GetResource(f.GetMessageTypes()[0]); got != nil {
			t.Errorf(`Got "%v", expected nil.`, got)
		}
	})
	t.Run("Nil", func(t *testing.T) {
		if got := GetResource(nil); got != nil {
			t.Errorf(`Got "%v", expected nil.`, got)
		}
	})
}

func TestGetResourceDefinition(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
		`)
		if got := GetResourceDefinitions(f); got != nil {
			t.Errorf("Got %v, expected nil.", got)
		}
	})
	t.Run("One", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			option (google.api.resource_definition) = {
				type: "library.googleapis.com/Book"
			};
		`)
		defs := GetResourceDefinitions(f)
		if got, want := len(defs), 1; got != want {
			t.Errorf("Got %d definitions, expected %d.", got, want)
		}
		if got, want := defs[0].GetType(), "library.googleapis.com/Book"; got != want {
			t.Errorf("Got %s for type, expected %s.", got, want)
		}
	})
	t.Run("Two", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			option (google.api.resource_definition) = {
				type: "library.googleapis.com/Book"
			};
			option (google.api.resource_definition) = {
				type: "library.googleapis.com/Author"
			};
		`)
		defs := GetResourceDefinitions(f)
		if got, want := len(defs), 2; got != want {
			t.Errorf("Got %d definitions, expected %d.", got, want)
		}
		if got, want := defs[0].GetType(), "library.googleapis.com/Book"; got != want {
			t.Errorf("Got %s for type, expected %s.", got, want)
		}
		if got, want := defs[1].GetType(), "library.googleapis.com/Author"; got != want {
			t.Errorf("Got %s for type, expected %s.", got, want)
		}
	})
}

func TestGetResourceReference(t *testing.T) {
	t.Run("Present", func(t *testing.T) {
		f := testutils.ParseProto3String(t, `
			import "google/api/resource.proto";
			message GetBookRequest {
				string name = 1 [(google.api.resource_reference) = {
					type: "library.googleapis.com/Book"
				}];
			}
		`)
		ref := GetResourceReference(f.GetMessageTypes()[0].GetFields()[0])
		if got, want := ref.GetType(), "library.googleapis.com/Book"; got != want {
			t.Errorf("Got %q, expected %q.", got, want)
		}
	})
	t.Run("Absent", func(t *testing.T) {
		f := testutils.ParseProto3String(t, "message GetBookRequest { string name = 1; }")
		if got := GetResourceReference(f.GetMessageTypes()[0].GetFields()[0]); got != nil {
			t.Errorf(`Got "%v", expected nil`, got)
		}
	})
}

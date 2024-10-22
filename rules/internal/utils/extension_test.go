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
	apb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
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
				t.Error(diff)
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
		{
			"Two",
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
				t.Error(diff)
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

func TestGetOperationInfoResponseType(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName     string
		ResponseType string
		valid        bool
	}{
		{"Valid", "WriteBookResponse", true},
		{"Invalid", "Foo", false},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			fd := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{ .ResponseType }}"
							metadata_type: "WriteBookMetadata"
						};
					}
				}
				message WriteBookRequest {}
				message WriteBookResponse {}
			`, test)

			typ := GetOperationResponseType(fd.GetServices()[0].GetMethods()[0])

			if validType := typ != nil; validType != test.valid {
				t.Fatalf("Expected valid(%v) response_type message", test.valid)
			}

			if !test.valid {
				return
			}

			if got, want := typ.GetName(), test.ResponseType; got != want {
				t.Errorf("Response type - got %q, want %q.", got, want)
			}
		})
	}
}

func TestGetOperationInfoMetadataType(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName     string
		MetadataType string
		valid        bool
	}{
		{"Valid", "WriteBookMetadata", true},
		{"Invalid", "Foo", false},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			fd := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";
				service Library {
					rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "WriteBookResponse"
							metadata_type: "{{ .MetadataType }}"
						};
					}
				}
				message WriteBookRequest {}
				message WriteBookMetadata {}
			`, test)

			typ := GetMetadataType(fd.GetServices()[0].GetMethods()[0])

			if validType := typ != nil; validType != test.valid {
				t.Fatalf("Expected valid(%v) metadata_type message", test.valid)
			}

			if !test.valid {
				return
			}

			if got, want := typ.GetName(), test.MetadataType; got != want {
				t.Errorf("Metadata type - got %q, want %q.", got, want)
			}
		})
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

func TestFindResource(t *testing.T) {
	files := testutils.ParseProtoStrings(t, map[string]string{
		"book.proto": `
			syntax = "proto3";
			package test;

			import "google/api/resource.proto";

			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};

				string name = 1;
			}
		`,
		"shelf.proto": `
			syntax = "proto3";
			package test;

			import "book.proto";
			import "google/api/resource.proto";

			message Shelf {
				option (google.api.resource) = {
					type: "library.googleapis.com/Shelf"
					pattern: "shelves/{shelf}"
				};

				string name = 1;

				repeated Book books = 2;
			}
		`,
	})

	for _, tst := range []struct {
		name, reference string
		notFound        bool
	}{
		{"local_reference", "library.googleapis.com/Shelf", false},
		{"imported_reference", "library.googleapis.com/Book", false},
		{"unresolvable", "foo.googleapis.com/Bar", true},
	} {
		t.Run(tst.name, func(t *testing.T) {
			got := FindResource(tst.reference, files["shelf.proto"])

			if tst.notFound && got != nil {
				t.Fatalf("Expected to not find the resource, but found %q", got.GetType())
			}

			if !tst.notFound && got == nil {
				t.Errorf("Got nil, expected %q", tst.reference)
			} else if !tst.notFound && got.GetType() != tst.reference {
				t.Errorf("Got %q, expected %q", got.GetType(), tst.reference)
			}
		})
	}
}

func TestFindResourceMessage(t *testing.T) {
	files := testutils.ParseProtoStrings(t, map[string]string{
		"book.proto": `
			syntax = "proto3";
			package test;

			import "google/api/resource.proto";

			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};

				string name = 1;
			}
		`,
		"shelf.proto": `
			syntax = "proto3";
			package test;

			import "book.proto";
			import "google/api/resource.proto";

			message Shelf {
				option (google.api.resource) = {
					type: "library.googleapis.com/Shelf"
					pattern: "shelves/{shelf}"
				};

				string name = 1;

				repeated Book books = 2;
			}
		`,
	})

	for _, tst := range []struct {
		name, reference, wantMsg string
		notFound                 bool
	}{
		{"local_reference", "library.googleapis.com/Shelf", "Shelf", false},
		{"imported_reference", "library.googleapis.com/Book", "Book", false},
		{"unresolvable", "foo.googleapis.com/Bar", "", true},
	} {
		t.Run(tst.name, func(t *testing.T) {
			got := FindResourceMessage(tst.reference, files["shelf.proto"])

			if tst.notFound && got != nil {
				t.Fatalf("Expected to not find the message, but found %q", got.GetName())
			}

			if !tst.notFound && got == nil {
				t.Errorf("Got nil, expected %q", tst.wantMsg)
			} else if !tst.notFound && got.GetName() != tst.wantMsg {
				t.Errorf("Got %q, expected %q", got.GetName(), tst.wantMsg)
			}
		})
	}
}

func TestSplitResourceTypeName(t *testing.T) {
	for _, tst := range []struct {
		name, input, service, typeName string
		ok                             bool
	}{
		{"Valid", "foo.googleapis.com/Foo", "foo.googleapis.com", "Foo", true},
		{"InvalidExtraSlashes", "foo.googleapis.com/Foo/Bar", "", "", false},
		{"InvalidNoService", "/Foo", "", "", false},
		{"InvalidNoTypeName", "foo.googleapis.com/", "", "", false},
	} {
		t.Run(tst.name, func(t *testing.T) {
			s, typ, ok := SplitResourceTypeName(tst.input)
			if ok != tst.ok {
				t.Fatalf("Expected %v for ok, but got %v", tst.ok, ok)
			}
			if diff := cmp.Diff(s, tst.service); diff != "" {
				t.Errorf("service: got(-),want(+):\n%s", diff)
			}
			if diff := cmp.Diff(typ, tst.typeName); diff != "" {
				t.Errorf("type name: got(-),want(+):\n%s", diff)
			}
		})
	}
}

func TestGetOutputOrLROResponseMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		RPCs string
		want string
	}{
		{"BookOutputType", `
			rpc CreateBook(CreateBookRequest) returns (Book) {};
		`, "Book"},
		{"BespokeOperationResource", `
			rpc CreateBook(CreateBookRequest) returns (Operation) {};
		`, "Operation"},
		{"LROBookResponse", `
			rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
				option (google.longrunning.operation_info) = {
					response_type: "Book"
				};
		};
		`, "Book"},
		{"LROMissingResponse", `
			rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
		};
		`, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				import "google/longrunning/operations.proto";
				import "google/protobuf/field_mask.proto";
				service Foo {
					{{.RPCs}}
				}

				// This is at the top to make it retrievable
				// by the test code.
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "books/{book}"
						singular: "book"
						plural: "books"
					};
				}

				message CreateBookRequest {
					// The parent resource where this book will be created.
					// Format: publishers/{publisher}
					string parent = 1;

					// The book to create.
					Book book = 2;
				}

				// bespoke operation message (not an LRO)
				message Operation {
				}
			`, test)
			method := file.GetServices()[0].GetMethods()[0]
			resp := GetResponseType(method)
			got := ""
			if resp != nil {
				got = resp.GetName()
			}
			if got != test.want {
				t.Errorf(
					"GetOutputOrLROResponseMessage got %q, want %q",
					got, test.want,
				)
			}
		})
	}
}

func TestFindResourceChildren(t *testing.T) {
	publisher := &apb.ResourceDescriptor{
		Type: "library.googleapis.com/Publisher",
		Pattern: []string{
			"publishers/{publisher}",
		},
	}
	shelf := &apb.ResourceDescriptor{
		Type: "library.googleapis.com/Shelf",
		Pattern: []string{
			"shelves/{shelf}",
		},
	}
	book := &apb.ResourceDescriptor{
		Type: "library.googleapis.com/Book",
		Pattern: []string{
			"publishers/{publisher}/books/{book}",
		},
	}
	edition := &apb.ResourceDescriptor{
		Type: "library.googleapis.com/Edition",
		Pattern: []string{
			"publishers/{publisher}/books/{book}/editions/{edition}",
		},
	}
	files := testutils.ParseProtoStrings(t, map[string]string{
		"book.proto": `
			syntax = "proto3";
			package test;

			import "google/api/resource.proto";

			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};

				string name = 1;
			}

			message Edition {
				option (google.api.resource) = {
					type: "library.googleapis.com/Edition"
					pattern: "publishers/{publisher}/books/{book}/editions/{edition}"
				};

				string name = 1;
			}
		`,
		"shelf.proto": `
			syntax = "proto3";
			package test;

			import "book.proto";
			import "google/api/resource.proto";

			message Shelf {
				option (google.api.resource) = {
					type: "library.googleapis.com/Shelf"
					pattern: "shelves/{shelf}"
				};

				string name = 1;

				repeated Book books = 2;
			}
		`,
	})

	for _, tst := range []struct {
		name   string
		parent *apb.ResourceDescriptor
		want   []*apb.ResourceDescriptor
	}{
		{"has_child_same_file", book, []*apb.ResourceDescriptor{edition}},
		{"has_child_other_file", publisher, []*apb.ResourceDescriptor{book, edition}},
		{"no_children", shelf, nil},
	} {
		t.Run(tst.name, func(t *testing.T) {
			got := FindResourceChildren(tst.parent, files["shelf.proto"])
			if diff := cmp.Diff(tst.want, got, cmp.Comparer(proto.Equal)); diff != "" {
				t.Errorf("got(-),want(+):\n%s", diff)
			}
		})
	}
}

func TestHasFieldInfo(t *testing.T) {
	testCases := []struct {
		name, FieldInfo string
		want            bool
	}{
		{
			name:      "HasFieldInfo",
			FieldInfo: "[(google.api.field_info).format = UUID4]",
			want:      true,
		},
		{
			name: "NoFieldInfo",
			want: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_info.proto";
			
			message CreateBookRequest {
				string foo = 1 {{.FieldInfo}};
			}
			`, tc)
			fd := file.FindMessage("CreateBookRequest").FindFieldByName("foo")
			if got := HasFieldInfo(fd); got != tc.want {
				t.Errorf("HasFieldInfo(%+v): expected %v, got %v", fd, tc.want, got)
			}
		})
	}
}

func TestGetFieldInfo(t *testing.T) {
	testCases := []struct {
		name, FieldInfo string
		want            *apb.FieldInfo
	}{
		{
			name:      "HasFieldInfo",
			FieldInfo: "[(google.api.field_info).format = UUID4]",
			want:      &apb.FieldInfo{Format: apb.FieldInfo_UUID4},
		},
		{
			name: "NoFieldInfo",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_info.proto";
			
			message CreateBookRequest {
				string foo = 1 {{.FieldInfo}};
			}
			`, tc)
			fd := file.FindMessage("CreateBookRequest").FindFieldByName("foo")
			got := GetFieldInfo(fd)
			if diff := cmp.Diff(got, tc.want, cmp.Comparer(proto.Equal)); diff != "" {
				t.Errorf("GetFieldInfo(%+v): got(-),want(+):\n%s", fd, diff)
			}
		})
	}
}

func TestHasFormat(t *testing.T) {
	testCases := []struct {
		name, Format string
		want         bool
	}{
		{
			name:   "HasFormat",
			Format: "format: UUID4",
			want:   true,
		},
		{
			name: "NoFormat",
			want: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_info.proto";
			
			message CreateBookRequest {
				string foo = 1 [(google.api.field_info) = {
					{{.Format}}
				}];
			}
			`, tc)
			fd := file.FindMessage("CreateBookRequest").FindFieldByName("foo")
			if got := HasFormat(fd); got != tc.want {
				t.Errorf("HasFormat(%+v): expected %v, got %v", fd, tc.want, got)
			}
		})
	}
}

func TestGetFormat(t *testing.T) {
	testCases := []struct {
		name, Format string
		want         apb.FieldInfo_Format
	}{
		{
			name:   "HasUUID4Format",
			Format: "format: UUID4",
			want:   apb.FieldInfo_UUID4,
		},
		{
			name: "NoFormat",
			want: apb.FieldInfo_FORMAT_UNSPECIFIED,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
			import "google/api/field_info.proto";
			
			message CreateBookRequest {
				string foo = 1 [(google.api.field_info) = {
					{{.Format}}
				}];
			}
			`, tc)
			fd := file.FindMessage("CreateBookRequest").FindFieldByName("foo")
			if got := GetFormat(fd); got != tc.want {
				t.Errorf("GetFormat(%+v): expected %v, got %v", fd, tc.want, got)
			}
		})
	}
}

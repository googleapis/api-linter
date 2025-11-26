package utils

import (
	"testing"

	"github.com/googleapis/api-linter/v2/rules/internal/testutils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestLintSingularStringField(t *testing.T) {
	for _, test := range []struct {
		testName  string
		FieldType string
		problems  testutils.Problems
	}{
		{"Valid", `string`, nil},
		{"Invalid", `int32`, testutils.Problems{{Suggestion: "string"}}},
		{"InvalidRepeated", `repeated string`, testutils.Problems{{Suggestion: "string"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Message {
					{{.FieldType}} foo = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			problems := LintSingularStringField(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintRequiredField(t *testing.T) {
	for _, test := range []struct {
		testName   string
		Annotation string
		problems   testutils.Problems
	}{
		{"Valid", `[(google.api.field_behavior) = REQUIRED]`, nil},
		{"Invalid", ``, testutils.Problems{{Message: "REQUIRED"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/field_behavior.proto";
				message Message {
					string foo = 1 {{.Annotation}};
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			problems := LintRequiredField(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintFieldResourceReference(t *testing.T) {
	for _, test := range []struct {
		testName   string
		Annotation string
		problems   testutils.Problems
	}{
		{"Valid", `[(google.api.resource_reference).type = "bar"]`, nil},
		{"Invalid", ``, testutils.Problems{{Message: "resource_reference"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Message {
					string foo = 1 {{.Annotation}};
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			problems := LintFieldResourceReference(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintNoHTTPBody(t *testing.T) {
	for _, test := range []struct {
		testName string
		Body     string
		problems testutils.Problems
	}{
		{"Valid", ``, nil},
		{"Invalid", `*`, testutils.Problems{{Message: "not have an HTTP body"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "/v1/{name=publishers/*/books/*}"
							body: "{{.Body}}"
						};
					}
				}
				message Book {}
				message GetBookRequest {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintNoHTTPBody(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintWildcardHTTPBody(t *testing.T) {
	for _, test := range []struct {
		testName string
		Body     string
		problems testutils.Problems
	}{
		{"Valid", `*`, nil},
		{"Invalid", ``, testutils.Problems{{Message: `use "*" as the HTTP body`}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc ArchiveBook(ArchiveBookRequest) returns (Book) {
						option (google.api.http) = {
							post: "/v1/{name=publishers/*/books/*}:archive"
							body: "{{.Body}}"
						};
					}
				}
				message Book {}
				message ArchiveBookRequest {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintWildcardHTTPBody(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintHTTPMethod(t *testing.T) {
	for _, test := range []struct {
		testName string
		Method   string
		problems testutils.Problems
	}{
		{"Valid", `get`, nil},
		{"Invalid", `delete`, testutils.Problems{{Message: `HTTP GET`}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							{{.Method}}: "/v1/{name=publishers/*/books/*}"
						};
					}
				}
				message Book {}
				message GetBookRequest {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintHTTPMethod("GET")(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintMethodHasMatchingRequestName(t *testing.T) {
	for _, test := range []struct {
		testName    string
		MessageName string
		problems    testutils.Problems
	}{
		{"Valid", "GetBookRequest", nil},
		{"Invalid", "AcquireBookRequest", testutils.Problems{{Suggestion: "GetBookRequest"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc GetBook({{.MessageName}}) returns (Book);
				}
				message Book {}
				message {{.MessageName}} {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintMethodHasMatchingRequestName(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintMethodHasMatchingResponseName(t *testing.T) {
	for _, test := range []struct {
		testName     string
		ResponseName string
		problems     testutils.Problems
	}{
		{"Valid", "GetBookResponse", nil},
		{"Invalid", "AcquireBookResponse", testutils.Problems{{Suggestion: "GetBookResponse"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc GetBook(GetBookRequest) returns ({{.ResponseName}});
				}
				message GetBookRequest {}
				message {{.ResponseName}} {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintMethodHasMatchingResponseName(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintMethodHasMatchingResponseNameLRO(t *testing.T) {
	for _, test := range []struct {
		testName    string
		MessageName string
		problems    testutils.Problems
	}{
		{"Valid", "GetBookResponse", nil},
		{"Invalid", "AcquireBookResponse", testutils.Problems{{Message: "GetBookResponse"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/longrunning/operations.proto";

				service Library {
					rpc GetBook(GetBookRequest) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{.MessageName}}"
							metadata_type: "OperationMetadata"
						};
					}
				}
				message GetBookRequest {}
				message {{.MessageName}} {}
				message OperationMetadata {}
			`, test)
			method := f.Services().Get(0).Methods().Get(0)
			problems := LintMethodHasMatchingResponseName(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintSingularField(t *testing.T) {
	for _, test := range []struct {
		testName string
		Label    string
		problems testutils.Problems
	}{
		{"Valid", "", nil},
		{"Invalid", "repeated", testutils.Problems{{Suggestion: "string"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Message {
					{{.Label}} string foo = 1;
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			problems := LintSingularField(field, protoreflect.StringKind, "string")
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintNotOneof(t *testing.T) {
	for _, test := range []struct {
		testName string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", `string foo = 1;`, nil},
		{"ValidProto3Optional", `optional string foo = 1;`, nil},
		{"Invalid", `oneof foo_oneof { string foo = 1; }`, testutils.Problems{{Message: "should not be a oneof"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Message {
					{{.Field}}
				}
			`, test)
			field := f.Messages().Get(0).Fields().Get(0)
			problems := LintNotOneof(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintPluralMethodName(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName       string
		prefix         string
		MethodName     string
		CollectionName string
		ResponseItems  string
		problems       testutils.Problems
	}{
		{
			testName:       "ValidBatchGetBooks",
			prefix:         "BatchGet",
			MethodName:     "BatchGetBooks",
			CollectionName: "books",
			ResponseItems:  "repeated Book books = 1;",
			problems:       testutils.Problems{},
		},
		{
			testName:       "ValidBatchGetMen",
			prefix:         "BatchGet",
			MethodName:     "BatchGetMen",
			CollectionName: "men",
			ResponseItems:  "repeated Other men = 1;",
			problems:       testutils.Problems{},
		},
		{
			testName:       "ValidBatchGetNames-NonMessageItems",
			prefix:         "BatchGet",
			MethodName:     "BatchGetNames",
			CollectionName: "names",
			ResponseItems:  "repeated string names = 1;",
			problems:       testutils.Problems{},
		},
		{
			testName:       "InvalidSingularBus",
			prefix:         "BatchGet",
			MethodName:     "BatchGetBus",
			CollectionName: "buses",
			ResponseItems:  "repeated Other buses = 1;",
			problems:       testutils.Problems{{Message: "Buses", Suggestion: "BatchGetBuses"}},
		},
		{
			testName:       "Invalid-SingularCorpPerson",
			prefix:         "BatchGet",
			MethodName:     "BatchGetCorpPerson",
			CollectionName: "corpPerson",
			ResponseItems:  "repeated Other corp_people = 1;",
			problems:       testutils.Problems{{Message: "CorpPeople", Suggestion: "BatchGetCorpPeople"}},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";

				service Test {
					rpc {{.MethodName}}({{.MethodName}}Request) returns ({{.MethodName}}Response) {
						option (google.api.http) = {
							get: "/v1/{parent=publishers/*}/{{.CollectionName}}:batchGet"
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

			problems := LintPluralMethodName(m, test.prefix)
			if diff := test.problems.SetDescriptor(m).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

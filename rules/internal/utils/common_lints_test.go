package utils

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestLintStringField(t *testing.T) {
	for _, test := range []struct {
		testName  string
		FieldType string
		problems  testutils.Problems
	}{
		{"Valid", `string`, nil},
		{"Invalid", `int32`, testutils.Problems{{Suggestion: "string"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Message {
					{{.FieldType}} foo = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			problems := LintStringField(field)
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
			field := f.GetMessageTypes()[0].GetFields()[0]
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
			field := f.GetMessageTypes()[0].GetFields()[0]
			problems := LintFieldResourceReference(field)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestLintParentFieldPresenceAndType(t *testing.T) {
	for _, test := range []struct {
		testName string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", `string parent = 1;`, nil},
		{"Missing", ``, testutils.Problems{{Message: "has no `parent` field"}}},
		{"WrongType", `bytes parent = 1;`, testutils.Problems{{Suggestion: "string"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message Message {
					{{.Field}}
					int32 other = 2;
				}
			`, test)
			message := f.GetMessageTypes()[0]
			var d desc.Descriptor = message
			if test.testName == "WrongType" {
				d = message.GetFields()[0]
			}
			problems := LintParentFieldPresenceAndType(message)
			if diff := test.problems.SetDescriptor(d).Diff(problems); diff != "" {
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
			method := f.GetServices()[0].GetMethods()[0]
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
			method := f.GetServices()[0].GetMethods()[0]
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
			method := f.GetServices()[0].GetMethods()[0]
			problems := LintHTTPMethod("GET")(method)
			if diff := test.problems.SetDescriptor(method).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

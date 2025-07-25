package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestParentRequired(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "ListBooksRequest", "parent", nil},
		{"InvalidName", "ListBooksRequest", "publisher", testutils.Problems{{Message: "no `parent` field"}}},
		{"Irrelevant", "EnumerateBooksRequest", "id", nil},
		{"IrrelevantAIP162", "ListBookRevisionsRequest", "name", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					string {{.FieldName}} = 1;
				}
			`, test)
			problems := requestParentRequired.Lint(f)
			if diff := test.problems.SetDescriptor(f.Messages().Get(0)).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}

	// Test the "top-level exception", which is more involved
	// than the other tests and therefore handled separately.
	for _, test := range []struct {
		testName string
		Package  string
	}{
		{"ValidTopLevel", ""},
		{"ValidTopLevelWithPackage", "package foo;"},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				{{.Package}}
				import "google/api/resource.proto";
				message ListBooksRequest {}
				message ListBooksResponse {
					repeated Book books = 1;
				}
				message Book {
					option (google.api.resource) = {
						pattern: "books/{book}"
					};
					string name = 1;
				}
			`, test)
			problems := requestParentRequired.Lint(f)
			if diff := (testutils.Problems{}).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

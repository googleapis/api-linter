package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestShowDeletedRequired(t *testing.T) {
	tests := []struct {
		name     string
		Message  string
		Method   string
		Field    string
		problems testutils.Problems
	}{
		{"Valid", "ListBooksRequest", "UndeleteBook", `bool show_deleted = 1;`, nil},
		{"Invalid", "ListBooksRequest", "UndeleteBook", ``, testutils.Problems{{Message: "show_deleted"}}},
		{"IrrelevantNoSoftDelete", "ListBooksRequest", "GetBook", ``, nil},
		{"IrrelevantMessage", "EnumerateBooksRequest", "UndeleteBook", "", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				service Library {
					rpc {{.Method}}({{.Method}}Request) returns (Book);
				}
				message Book {}
				message {{.Message}} {
					{{.Field}}
				}
				message ListBooksResponse {
					repeated Book books = 1;
				}
				message {{.Method}}Request {}
			`, test)
			problems := requestShowDeletedRequired.Lint(f)
			if diff := test.problems.SetDescriptor(f.GetMessageTypes()[1]).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

// Regression test for https://github.com/googleapis/api-linter/issues/854.
func TestRequestShowDeletedRequired_NonMessageType(t *testing.T) {
	f := testutils.ParseProto3Tmpl(t, `
		message ListBooksRequest {}
		message ListBooksResponse {
			repeated string books = 1;
		}
	`, nil)
	problems := requestShowDeletedRequired.Lint(f)
	if diff := (testutils.Problems{}).SetDescriptor(f.GetMessageTypes()[0]).Diff(problems); diff != "" {
		t.Error(diff)
	}
}

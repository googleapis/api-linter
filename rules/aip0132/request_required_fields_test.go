package aip0132

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestRequestRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		MessageName string
		FieldName   string
		problems    testutils.Problems
	}{
		{"Valid", "ListBooksRequest", "parent", nil},
		{"InvalidName", "ListBooksRequest", "publisher", testutils.Problems{{Message: "no `parent` field"}}},
		{"Irrelevant", "EnumerateBooksRequest", "id", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				message {{.MessageName}} {
					string {{.FieldName}} = 1;
				}
			`, test)
			problems := requestRequiredFields.Lint(f)
			if diff := test.problems.SetDescriptor(f.GetMessageTypes()[0]).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

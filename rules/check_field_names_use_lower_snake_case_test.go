package rules

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/rules/testdata"
)

func TestFieldNamesUseLowerSnakeCaseRule(t *testing.T) {
	tmpl := testdata.MustCreateTemplate(`
	syntax = "proto2";
	message Foo {
	  optional string {{.FieldName}} = 1;
	}`)

	tests := []struct {
		FieldName  string
		numProblem int
		suggestion string
		startLine  int
	}{
		{"good_field_name", 0, "", -1},
		{"BadFieldName", 1, "bad_field_name", 4},
		{"badFieldName", 1, "bad_field_name", 4},
		{"Bad_Field_Name", 1, "bad_field_name", 4},
		{"bad_Field_Name", 1, "bad_field_name", 4},
		{"badField_Name", 1, "bad_field_name", 4},
	}

	rule := checkFieldNamesUseLowerSnakeCase()

	for _, test := range tests {
		req := testdata.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := fmt.Sprintf("Check field name `%s`", test.FieldName)
		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp), test.numProblem; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp) > 0 {
			if got, want := resp[0].Suggestion, test.suggestion; got != want {
				t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, got, want)
			}
			if got, want := resp[0].Location.Start.Line, test.startLine; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

package rules

import (
	"fmt"
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/testdata"
)

func TestFieldNamesUseLowerSnakeCaseRule(t *testing.T) {
	tmpl := testdata.MustCreateTemplateWithDedent(`
	syntax = "proto2";
	message Foo {
	  optional string {{.FieldName}} = 1;
	}`)

	wantPosition := lint.Position{Line: 3, Column: 3}
	tests := []struct {
		FieldName  string
		numProblem int
		suggestion string
		start      lint.Position
	}{
		{"good_field_name", 0, "", lint.Position{}},
		{"BadFieldName", 1, "bad_field_name", wantPosition},
		{"badFieldName", 1, "bad_field_name", wantPosition},
		{"Bad_Field_Name", 1, "bad_field_name", wantPosition},
		{"bad_Field_Name", 1, "bad_field_name", wantPosition},
		{"badField_Name", 1, "bad_field_name", wantPosition},
	}

	rule := checkFieldNamesUseLowerSnakeCase()

	for _, test := range tests {
		req := testdata.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := fmt.Sprintf("Check field name `%s`", test.FieldName)
		resp, err := rule.Lint(req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}

		if got, want := len(resp.Problems), test.numProblem; got != want {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, got, want)
		}

		if len(resp.Problems) > 0 {
			if got, want := resp.Problems[0].Suggestion, test.suggestion; got != want {
				t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, got, want)
			}
			if got, want := resp.Problems[0].Location.Start, test.start; got != want {
				t.Errorf("%s: got location starting with %d, but want %d", errPrefix, got, want)
			}
		}
	}
}

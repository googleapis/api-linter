package rules

import (
	"fmt"
	"testing"

	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules/testdata"
)

func TestFieldNamesUseLowerSnakeCase(t *testing.T) {
	tmpl := testdata.MustCreateTemplate(`
	syntax = "proto2";
	message Foo {
	  optional string {{.FieldName}} = 1;
	}`)

	tests := []struct {
		FieldName  string
		numProblem int
		suggestion string
	}{
		{"good_field_name", 0, ""},
		{"BadFieldName", 1, "bad_field_name"},
		{"badFieldName", 1, "bad_field_name"},
		{"Bad_Field_Name", 1, "bad_field_name"},
		{"bad_Field_Name", 1, "bad_field_name"},
		{"badField_Name", 1, "bad_field_name"},
	}

	rules, err := lint.NewRules(checkNamingFormats())
	if err != nil {
		t.Errorf("lint.NewRules return error %v", err)
	}

	for _, test := range tests {
		req := testdata.MustCreateRequestFromTemplate(tmpl, test)

		errPrefix := fmt.Sprintf("Check field name `%s`", test.FieldName)
		resp, err := lint.Run(rules, req)
		if err != nil {
			t.Errorf("%s: lint.Run return error %v", errPrefix, err)
		}
		if len(resp.Problems) != test.numProblem {
			t.Errorf("%s: got %d problems, but want %d", errPrefix, len(resp.Problems), test.numProblem)
		}
		if len(resp.Problems) > 0 && resp.Problems[0].Suggestion != test.suggestion {
			t.Errorf("%s: got suggestion '%s', but want '%s'", errPrefix, resp.Problems[0].Suggestion, test.suggestion)
		}
	}
}

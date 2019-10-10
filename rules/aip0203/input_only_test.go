package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestInputOnly(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "valid case with INPUT_ONLY field behavior annotation",
			comment:  "input_only",
			field:    "string secret = 1 [(google.api.field_behavior) = INPUT_ONLY];",
			problems: nil,
		},
		{
			name:    "input_only",
			comment: "input_only",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "inputonly",
			comment: "inputonly",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "@inputonly",
			comment: "@inputonly",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "@input_only",
			comment: "@input_only",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "INPUT_ONLY",
			comment: "INPUT_ONLY",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "inputOnly",
			comment: "inputOnly",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "@inputOnly",
			comment: "@inputonly",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "@INPUT_ONLY",
			comment: "@INPUT_ONLY",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "input_only_free_text",
			comment: "This field is input only",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "!inputOnly",
			comment: "!inputOnly",
			field:   "string secret = 1;",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
				import "google/api/field_behavior.proto";

				message Book {
					// Secrets to be stored in the book
					// {{.Comment}}
					{{.Field}}
				}
				`
			file := testutils.ParseProto3Tmpl(t, template, struct {
				Comment string
				Field   string
			}{test.comment, test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := inputOnly.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

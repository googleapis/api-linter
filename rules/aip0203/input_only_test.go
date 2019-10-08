package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

func TestInputOnly(t *testing.T) {
	testCases := []struct {
		name       string
		annotation string
		problems   testutils.Problems
	}{
		{
			name:       "input_only",
			annotation: "input_only",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "inputonly",
			annotation: "inputonly",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "@inputonly",
			annotation: "@inputonly",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "@input_only",
			annotation: "@input_only",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "INPUT_ONLY",
			annotation: "INPUT_ONLY",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "inputOnly",
			annotation: "inputOnly",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "@inputOnly",
			annotation: "@inputonly",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "@INPUT_ONLY",
			annotation: "@INPUT_ONLY",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "input_only_free_text",
			annotation: "This field is input only",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:       "!inputOnly",
			annotation: "!inputOnly",
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
				message Book {
					// Secrets to be stored in the book
					// {{.Annotation}}
					string secrect = 1;
				}
				`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Annotation string }{test.annotation})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := inputOnly.LintField(f)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

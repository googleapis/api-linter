package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

var title = "string title = 1;"
var titleWithRequiredBehavior = "string title = 1 [(google.api.field_behavior) = REQUIRED];"

func TestRequired(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "Valid",
			comment:  "Required",
			field:    titleWithRequiredBehavior,
			problems: nil,
		},
		{
			name:     "Valid",
			comment:  "@required",
			field:    titleWithRequiredBehavior,
			problems: nil,
		},
		{
			name:    "Invalid-required",
			comment: "required",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-Required",
			comment: "Required",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@required",
			comment: "@required",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@Required",
			comment: "@Required",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-REQUIRED",
			comment: "REQUIRED",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@REQUIRED",
			comment: "@REQUIRED",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-required_free_text",
			comment: "This field is required",
			field:   title,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-!required",
			comment: "!required",
			field:   title,
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
	// Title of the book
	// {{.Comment}}
	{{.Field}}
}`
			file := testutils.ParseProto3Tmpl(t, template,
				struct {
					Comment string
					Field   string
				}{test.comment, test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := required.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

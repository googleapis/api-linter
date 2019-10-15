package aip0203

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
)

var titleField = "string title = 1;"
var titleWithOptionalBehavior = "string title = 1 [(google.api.field_behavior) = OPTIONAL];"

func TestOptional(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "Valid",
			comment:  "Optional",
			field:    titleWithOptionalBehavior,
			problems: nil,
		},
		{
			name:     "Valid",
			comment:  "@optional",
			field:    titleWithOptionalBehavior,
			problems: nil,
		},
		{
			name:    "Invalid-optional",
			comment: "optional",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-Optional",
			comment: "Optional",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@optional",
			comment: "@optional",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@Optional",
			comment: "@Optional",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-OPTIONAL",
			comment: "OPTIONAL",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-@OPTIONAL",
			comment: "@OPTIONAL",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-optional_free_text",
			comment: "This field is optional",
			field:   titleField,
			problems: testutils.Problems{{
				Message: "google.api.field_behavior",
			}},
		},
		{
			name:    "Invalid-!optional",
			comment: "!optional",
			field:   titleField,
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
			problems := optional.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestOptionalBehaviorConflict(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		problems testutils.Problems
	}{
		{
			name:     "Valid",
			field:    "string title = 1 [(google.api.field_behavior) = OPTIONAL];",
			problems: nil,
		},
		{
			name: "Valid",
			field: `
string title = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];`,
			problems: nil,
		},
		{
			name: "Invalid-optional-conflict",
			field: `
string title = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY,
	(google.api.field_behavior) = OPTIONAL];`,
			problems: testutils.Problems{{
				Message: "Field behavior `(google.api.field_behavior) = OPTIONAL` shouldn't be used together with other field behaviors",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
import "google/api/field_behavior.proto";
message Book {
	// Title of the book
	{{.Field}}
}`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Field string }{test.field})
			f := file.GetMessageTypes()[0].GetFields()[0]
			problems := optionalBehaviorConflict.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestOptionalBehaviorConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		problems testutils.Problems
	}{
		{
			name: "Valid-NoneOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3;

string author = 4;`,
			problems: nil,
		},
		{
			name: "Valid-AllOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3 [(google.api.field_behavior) = OPTIONAL];

string author = 4 [(google.api.field_behavior) = OPTIONAL];`,
			problems: nil,
		},
		{
			name: "Invalid-PartialOptional",
			field: `
string name = 1 [
	(google.api.field_behavior) = IMMUTABLE,
	(google.api.field_behavior) = OUTPUT_ONLY];

string title = 2 [(google.api.field_behavior) = REQUIRED];

string summary = 3 [(google.api.field_behavior) = OPTIONAL];

string author = 4;`,
			problems: testutils.Problems{{
				Message: "Within a single message, either all optional fields should be indicated, or none of them should be.",
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			template := `
import "google/api/field_behavior.proto";
message Book {
	// Title of the book
	{{.Field}}
}`
			file := testutils.ParseProto3Tmpl(t, template, struct{ Field string }{test.field})
			// author field in the test will get the warning
			f := file.GetMessageTypes()[0].GetFields()[3]
			problems := optionalBehaviorConsistency.Lint(file)
			if diff := test.problems.SetDescriptor(f).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// This rule inspects the leading comments of each field
// and if anything looks similar to "Optional.", it throws
// a problem.
//
// Examples:
// Incorrect code for this rule:
//
//	message Book {
//	// The title of the book.
//	// @Optional
//	string title = 1;
//	}
//
// or
//
//	message Book {
//	// Optional. The title of the book.
//	string title = 1;
//	}
//
//
// Correct code for this rule:
//
//	message Book {
//		// The title of the book.
//		string title = 1 [(google.api.field_behavior) = OPTIONAL];
//	}
var optional = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0203", "optional"),
	URI:    "http://aip.dev/203#optional",
	OnlyIf: withoutOptionalFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, optionalRegexp, "OPTIONAL")
	},
}

var optionalBehaviorConflict = &lint.FieldRule{
	Name: lint.NewRuleName("core", "0203", "optional-conflict"),
	URI:  "http://aip.dev/203#optional",
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return !withoutOptionalFieldBehavior(f)
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// APIs may use the OPTIONAL value to describe a field which doesn't use
		// REQUIRED, IMMUTABLE, INPUT_ONLY or OUTPUT_ONLY. If a field is described
		// as optional, it can't be others.
		if len(utils.GetFieldBehavior(f)) > 1 {
			return []lint.Problem{{
				Message:    "Field behavior `(google.api.field_behavior) = OPTIONAL` shouldn't be used together with other field behaviors.",
				Descriptor: f,
			}}
		}
		return nil
	},
}

// If a field is described as optional, ensure that every optional field on the
// message has this indicator.
var optionalBehaviorConsistency = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0203", "optional-consistency"),
	URI:    "https://aip.dev/203#optional",
	OnlyIf: hasOptionalFieldBehavior,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {

		for _, f := range m.GetFields() {
			if utils.GetFieldBehavior(f) == nil {
				problems = append(problems, lint.Problem{
					Message:    "Within a single message, either all optional fields should be indicated, or none of them should be.",
					Descriptor: f,
				})
			}
		}
		return
	},
}

var optionalRegexp = regexp.MustCompile("(?i).*optional.*")

func withoutOptionalFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_OPTIONAL {
			return false
		}
	}
	return true
}

func hasOptionalFieldBehavior(m *desc.MessageDescriptor) bool {
	for _, f := range m.GetFields() {
		for _, v := range utils.GetFieldBehavior(f) {
			if v == annotations.FieldBehavior_OPTIONAL {
				return true
			}
		}
	}
	return false
}

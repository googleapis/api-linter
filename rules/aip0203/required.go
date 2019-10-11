package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// This rule inspects the leading comments of each field
// and if anything looks similar to "Required.", it throws
// a problem.
//
// Examples:
// Incorrect code for this rule:
//
//	message Book {
//	// The title of the book.
//	// @Required
//	string title = 1;
//	}
//
// or
//
//	message Book {
//	// Required. The title of the book.
//	string title = 1;
//	}
//
//
// Correct code for this rule:
//
//	message Book {
//		// The title of the book.
//		string title = 1 [(google.api.field_behavior) = REQUIRED];
//	}
var required = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0203", "required"),
	URI:    "http://api.dev/203#guidance",
	OnlyIf: withoutRequiredRegexpFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, requiredRegexp, "REQUIRED")
	},
}

var requiredRegexp = regexp.MustCompile("(?i).*required.*")

func withoutRequiredRegexpFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_REQUIRED {
			return false
		}
	}
	return true
}

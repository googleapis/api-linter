package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// This rule inspects the leading comments of each field
// and if anything looks similar to "Input only.", it throws
// a problem.
//
// Examples:
// Incorrect code for this rule:
//
// message Book {
//   // Secrets to be stored in the book
//   // @InputOnly
//   string secret = 1;
// }
//
// or
//
// message Book {
//   // Input only. Secret to be stored in the book.
//   string secret = 1;
// }
//
// Correct code for this rule:
//
// message Book {
//   // Secret to be stored in the book.
//   string secret = 1 [(google.api.field_behavior) = INPUT_ONLY];
// }
var inputOnly = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0203", "input-only"),
	URI:    "http://api.dev/203#guidance",
	OnlyIf: withoutInputOnlyFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, inputOnlyRegexp, "INPUT_ONLY")
	},
}

var inputOnlyRegexp = regexp.MustCompile("(?i).*input.?only.*")

func withoutInputOnlyFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_INPUT_ONLY {
			return false
		}
	}
	return true
}

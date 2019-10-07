package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
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
	Name: lint.NewRuleName("core", "0203", "input-only"),
	URI:  "http://api.dev/203#guidance",
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkAnnotation(f, inputOnlyRegexp)
	},
}

var inputOnlyRegexp = regexp.MustCompile("(?i).*input.?only.*")

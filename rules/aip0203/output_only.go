package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// This rule inspects the leading comments of each field
// and if anything looks similar to "Output Only", it throws
// a problem.
//
// Examples:
// Incorrect code for this rule:
//
//	message Book {
//	// A generated URI for this book.
//	// @OutputOnly
//	string generated_uri = 1;
//	}
//
// or
//
//	message Book {
//	// Output only. A generated URI for this book.
//	string generated_uri = 1;
//	}
//
//
// Correct code for this rule:
//
//	message Book {
//	// A generated URI for this book.
//	string generated_uri = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
//	}
var outputOnly = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0203", "output-only"),
	URI:    "http://api.dev/203#guidance",
	OnlyIf: withoutOutputOnlyFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, outputOnlyRegexp, "OUTPUT_ONLY")
	},
}

var outputOnlyRegexp = regexp.MustCompile("(?i).*output.?only.*")

func withoutOutputOnlyFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_OUTPUT_ONLY {
			return false
		}
	}
	return true
}

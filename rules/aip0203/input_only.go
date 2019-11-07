package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var inputOnly = &lint.FieldRule{
	Name:   lint.NewRuleName(203, "input-only"),
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

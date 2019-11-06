package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var outputOnly = &lint.FieldRule{
	Name:   lint.NewRuleName(203, "output-only"),
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

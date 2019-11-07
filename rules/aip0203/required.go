package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var required = &lint.FieldRule{
	Name:   lint.NewRuleName(203, "required"),
	OnlyIf: withoutRequiredFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, requiredRegexp, "REQUIRED")
	},
}

var requiredRegexp = regexp.MustCompile("(?i).*required.*")

func withoutRequiredFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_REQUIRED {
			return false
		}
	}
	return true
}

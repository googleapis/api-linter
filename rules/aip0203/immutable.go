package aip0203

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
)

var immutable = &lint.FieldRule{
	Name:   lint.NewRuleName("core", "0203", "immutable"),
	OnlyIf: withoutImmutableFieldBehavior,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return checkLeadingComments(f, immutableRegexp, "IMMUTABLE")
	},
}

var immutableRegexp = regexp.MustCompile("(?i).*immutable.*")

func withoutImmutableFieldBehavior(f *desc.FieldDescriptor) bool {
	for _, v := range utils.GetFieldBehavior(f) {
		if v == annotations.FieldBehavior_IMMUTABLE {
			return false
		}
	}
	return true
}

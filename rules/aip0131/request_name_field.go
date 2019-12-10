package aip0131

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// Get request should have a string name field.
var requestNameField = &lint.FieldRule{
	Name: lint.NewRuleName(131, "request-name-field"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isGetRequestMessage(f.GetOwner()) && f.GetName() == "name"
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.GetType() != builder.FieldTypeString().GetType() {
			return []lint.Problem{{
				Message:    "`name` field on Get RPCs should be a string",
				Descriptor: f,
				Location:   locations.FieldType(f),
				Suggestion: "string",
			}}
		}

		return nil
	},
}

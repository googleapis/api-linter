// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0123

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceReferenceType = &lint.FieldRule{
	Name: lint.NewRuleName(123, "resource-reference-type"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.GetResourceReference(f) != nil
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if utils.GetTypeName(f) != "string" {
			// We assume that the likely mistake is probably that the annotation
			// is wrong (and should not be there), and not that the type is wrong,
			// because this is what we have observed in real life.
			return []lint.Problem{{
				Message:    "Resource references should only be applied to strings.",
				Descriptor: f,
				Location:   locations.FieldResourceReference(f),
				Suggestion: "",
			}}
		}
		return nil
	},
}

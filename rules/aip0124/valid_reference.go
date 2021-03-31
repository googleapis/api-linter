// Copyright 2020 Google LLC
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

package aip0124

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var validReference = &lint.FieldRule{
	Name:   lint.NewRuleName(124, "valid-reference"),
	OnlyIf: isUnknownType,
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// Get the type we are checking for.
		ref := utils.GetResourceReference(f)
		urt := ref.GetType()
		if urt == "" {
			urt = ref.GetChildType()
		}

		// Iterate over each dependency file and check for a matching resource.
		for _, file := range utils.GetAllDependencies(f.GetFile()) {
			// Most resources will be messages. If we find a message with a
			// resource annotation matching our universal resource type, we are done.
			for _, message := range file.GetMessageTypes() {
				if res := utils.GetResource(message); res != nil {
					if res.GetType() == urt {
						return nil
					}
				}
			}

			// Some resources are defined as file annotations. Check for these too.
			for _, rd := range utils.GetResourceDefinitions(file) {
				if rd.GetType() == urt {
					return nil
				}
			}
		}

		// We could not find a resource with that type. Return a problem.
		return []lint.Problem{{
			Message:    fmt.Sprintf("Could not find resource of type %q", urt),
			Descriptor: f,
			Location:   locations.FieldResourceReference(f),
		}}
	},
}

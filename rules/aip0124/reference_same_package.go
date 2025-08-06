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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var referenceSamePackage = &lint.FieldRule{
	Name:   lint.NewRuleName(124, "reference-same-package"),
	OnlyIf: isUnknownType,
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		// Get the type we are checking for.
		ref := utils.GetResourceReference(f)
		urt := ref.GetType()
		if urt == "" {
			urt = ref.GetChildType()
		}

		// Iterate over each dependency file and check for a matching resource.
		for _, file := range getNonPkgDependencies(f.ParentFile(), f.ParentFile().Package()) {
			// If we find a message with a resource annotation matching our universal
			// resource type, then it is in the wrong package.
			for i := 0; i < file.Messages().Len(); i++ {
				message := file.Messages().Get(i)
				if res := utils.GetResource(message); res != nil && res.GetType() == urt {
					return []lint.Problem{{
						Message:    fmt.Sprintf("Resource type %q should be declared in the same package as it is referenced.", urt),
						Descriptor: f,
						Location:   locations.FieldResourceReference(f),
					}}
				}
			}

			// Some resources are defined as file annotations. Check for these too.
			for _, rd := range utils.GetResourceDefinitions(file) {
				if rd.GetType() == urt {
					return []lint.Problem{{
						Message:    fmt.Sprintf("Resource type %q should be declared in the same package as it is referenced.", urt),
						Descriptor: f,
						Location:   locations.FieldResourceReference(f),
					}}
				}
			}
		}

		return nil
	},
}

// getNonPkgDependencies returns dependencies in other packages.
func getNonPkgDependencies(file protoreflect.FileDescriptor, pkg protoreflect.FullName) map[string]protoreflect.FileDescriptor {
	answer := map[string]protoreflect.FileDescriptor{}
	for name, dep := range utils.GetAllDependencies(file) {
		if dep.Package() != pkg {
			answer[name] = dep
		}
	}
	return answer
}

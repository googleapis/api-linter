// Copyright 2022 Google LLC
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
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceDefinitionTypeName = &lint.FileRule{
	Name:   lint.NewRuleName(123, "resource-definition-type-name"),
	OnlyIf: hasResourceDefinitionAnnotation,
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		var problems []lint.Problem
		resources := utils.GetResourceDefinitions(f)
		for ndx, resource := range resources {
			if strings.Count(resource.GetType(), "/") != 1 {
				problems = append(problems, lint.Problem{
					Message:    "Resource type names must be of the form {Service Name}/{Type}.",
					Descriptor: f,
					Location:   locations.FileResourceDefinition(f, ndx),
				})
			}
		}

		return problems
	},
}

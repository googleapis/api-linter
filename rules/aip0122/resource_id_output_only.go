// Copyright 2019 Google LLC
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

package aip0122

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/stoewer/go-strcase"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceIdOutputOnly = &lint.FieldRule{
	Name: lint.NewRuleName(122, "resource-id-output-only"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		var idName string
		p := f.Parent().(protoreflect.MessageDescriptor)

		// Build an expected ID field name based on the Resource `singular`
		// field or by parsing the `type`.
		isRes := utils.IsResource(p)
		if isRes {
			res := utils.GetResource(p)
			idName = res.GetSingular()
			if idName == "" {
				if _, t, ok := utils.SplitResourceTypeName(res.GetType()); ok {
					idName = strcase.SnakeCase(t)
				}
			}
			idName += "_id"
		}

		isId := f.Name() == "uid" || f.Name() == protoreflect.Name(idName)
		return isRes && isId
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		behaviors := utils.GetFieldBehavior(f)
		if !behaviors.Contains(annotations.FieldBehavior_OUTPUT_ONLY.String()) {
			return []lint.Problem{
				{
					Message:    "Resource ID fields must have field_behavior OUTPUT_ONLY",
					Descriptor: f,
				},
			}
		}
		return nil
	},
}

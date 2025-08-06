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

package aip0131

import (
	"fmt"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var requestNameReferenceType = &lint.FieldRule{
	Name: lint.NewRuleName(131, "request-name-reference-type"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		if m, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return utils.IsGetRequestMessage(m) && f.Name() == "name" && utils.GetResourceReference(f) != nil
		}
		return false
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if ref := utils.GetResourceReference(f); ref.GetType() == "" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("The `%s` field `google.api.resource_reference` annotation should be a direct `type` reference.", f.Name()),
				Descriptor: f,
			}}
		}
		return nil
	},
}

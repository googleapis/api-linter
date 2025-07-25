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

package aip0132

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var requestParentValidReference = &lint.FieldRule{
	Name: lint.NewRuleName(132, "request-parent-valid-reference"),
	OnlyIf: func(f protoreflect.FieldDescriptor) bool {
		ref := utils.GetResourceReference(f)
		if m, ok := f.Parent().(protoreflect.MessageDescriptor); ok {
			return utils.IsListRequestMessage(m) && f.Name() == "parent" && ref != nil && ref.GetType() != ""
		}
		return false
	},
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		p := f.Parent()
		msg := p.(protoreflect.MessageDescriptor)
		res := utils.GetResourceReference(f).GetType()

		response := utils.FindMessage(f.ParentFile(), strings.Replace(string(msg.Name()), "Request", "Response", 1))
		if response == nil {
			return nil
		}

		for i := 0; i < response.Fields().Len(); i++ {
			field := response.Fields().Get(i)
			typ := field.Message()
			if !field.IsList() && typ == nil {
				continue
			}

			if r := utils.GetResource(typ); r != nil && r.GetType() == res {
				return []lint.Problem{{
					Message:    fmt.Sprintf("The `google.api.resource_reference` on `%s` field should reference the parent(s) of `%s`.", f.Name(), res),
					Descriptor: f,
					Location:   locations.FieldResourceReference(f),
				}}
			}
		}

		return nil
	},
}
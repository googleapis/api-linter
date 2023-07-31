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
	"github.com/jhump/protoreflect/desc"
)

var requestParentValidReference = &lint.FieldRule{
	Name: lint.NewRuleName(132, "request-parent-valid-reference"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		ref := utils.GetResourceReference(f)
		return utils.IsListRequestMessage(f.GetOwner()) && f.GetName() == "parent" && ref != nil && ref.GetType() != ""
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		p := f.GetParent()
		msg := p.(*desc.MessageDescriptor)
		res := utils.GetResourceReference(f).GetType()

		response := utils.FindMessage(f.GetFile(), strings.Replace(msg.GetName(), "Request", "Response", 1))
		if response == nil {
			return nil
		}

		for _, field := range response.GetFields() {
			typ := field.GetMessageType()
			if !field.IsRepeated() && typ == nil {
				continue
			}

			if r := utils.GetResource(typ); r != nil && r.GetType() == res {
				return []lint.Problem{{
					Message:    fmt.Sprintf("The `google.api.resource_reference` on `%s` field should reference the parent(s) of `%s`.", f.GetName(), res),
					Descriptor: f,
					Location:   locations.FieldResourceReference(f),
				}}
			}
		}

		return nil
	},
}

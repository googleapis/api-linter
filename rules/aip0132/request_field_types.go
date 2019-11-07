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

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var knownFields = stringset.New("filter", "order_by")

// List methods should not have unrecognized fields.
var requestFieldTypes = &lint.FieldRule{
	Name: lint.NewRuleName(132, "request-field-types"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isListRequestMessage(f.GetOwner()) && knownFields.Contains(f.GetName())
	},
	LintField: func(f *desc.FieldDescriptor) (problems []lint.Problem) {
		// Establish that the field being checked is a string.
		if utils.GetScalarTypeName(f) != "string" {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Field %q should be a string.", f.GetName()),
				Descriptor: f,
			}}
		}
		return
	},
}

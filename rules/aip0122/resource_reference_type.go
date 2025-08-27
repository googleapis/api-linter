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

package aip0122

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var resourceReferenceType = &lint.FieldRule{
	Name: lint.NewRuleName(122, "resource-reference-type"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if utils.GetResourceReference(f) != nil && utils.GetTypeName(f) != "string" {
			return []lint.Problem{{
				Message:    "The resource_reference annotation should only be used on string fields.",
				Descriptor: f,
				Location:   locations.FieldResourceReference(f),
			}}
		}
		return nil
	},
}

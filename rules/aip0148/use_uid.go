// Copyright 2023 Google LLC
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

package aip0148

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

const (
	idStr  = "id"
	uidStr = "uid"
)

var useUid = &lint.FieldRule{
	Name: lint.NewRuleName(148, "use-uid"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.IsResource(f.GetOwner())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if f.GetName() == idStr {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Use %s instead of %s.", uidStr, idStr),
				Descriptor: f,
				Location:   locations.DescriptorName(f),
				Suggestion: uidStr,
			}}
		}
		return nil
	},
}

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

package aip0203

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	fpb "google.golang.org/genproto/googleapis/api/annotations"
)

var resourceNameIdentifier = &lint.MessageRule{
	Name:   lint.NewRuleName(203, "resource-name-identifier"),
	OnlyIf: utils.IsResource,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		f := m.FindFieldByName(utils.GetResourceNameField(utils.GetResource(m)))
		fb := utils.GetFieldBehavior(f)
		if len(fb) == 0 || !fb.Contains(fpb.FieldBehavior_IDENTIFIER.String()) {
			return []lint.Problem{{
				Message:    "resource name field must have field_behavior IDENTIFIER",
				Descriptor: f,
				Location:   locations.FieldOption(f, fpb.E_FieldBehavior),
			}}
		}

		return nil
	},
}

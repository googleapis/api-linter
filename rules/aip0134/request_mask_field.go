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

package aip0134

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var requestMaskField = &lint.MessageRule{
	Name:   lint.NewRuleName("core", "0134", "request-mask-field"),
	URI:    "https://aip.dev/134",
	OnlyIf: isUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		updateMask := m.FindFieldByName("update_mask")
		if updateMask == nil {
			return []lint.Problem{{
				Message:    "Update methods should have an `update_mask` field.",
				Descriptor: m,
			}}
		}
		if t := updateMask.GetMessageType(); t == nil || t.GetFullyQualifiedName() != "google.protobuf.FieldMask" {
			return []lint.Problem{{
				Message:    "The `update_mask` field should be a google.protobuf.FieldMask.",
				Descriptor: updateMask,
			}}
		}
		return nil
	},
}

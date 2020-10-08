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
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

// The resource itself is not included here, but also permitted.
// This is covered in code in the rule itself.
var respAllowedFields = stringset.New(
	"next_page_token",       // AIP-158
	"total_size",            // AIP-132
	"unreachable",           // AIP-217
	"unreachable_locations", // Wrong, but a separate AIP-217 rule catches it.
)

var responseUnknownFields = &lint.FieldRule{
	Name: lint.NewRuleName(132, "response-unknown-fields"),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return isListResponseMessage(f.GetOwner())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		// A repeated variant of the resource should be permitted.
		msgName := f.GetOwner().GetName()
		resource := strcase.SnakeCase(listRespMessageRegexp.FindStringSubmatch(msgName)[1])
		if strings.HasSuffix(resource, "_revisions") {
			// This is an AIP-162 ListFooRevisions response, which is subtly
			// different from an AIP-132 List response. We need to modify the RPC
			// name to what the AIP-132 List response would be in order to permit
			// the resource field properly.
			resource = utils.ToPlural(strings.TrimSuffix(resource, "_revisions"))
		}
		if f.GetName() == resource {
			return nil
		}

		// It is not the resource field; check it against the whitelist.
		if !respAllowedFields.Contains(f.GetName()) {
			return []lint.Problem{{
				Message:    "List responses should only contain fields explicitly described in AIPs.",
				Descriptor: f,
			}}
		}
		return nil
	},
}

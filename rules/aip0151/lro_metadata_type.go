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

package aip0151

import (
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/locations"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var lroMetadata = &lint.MethodRule{
	Name:   lint.NewRuleName(151, "lro-metadata-type"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		lro := utils.GetOperationInfo(m)

		// Ensure the response type is set.
		if lro.GetMetadataType() == "" {
			return []lint.Problem{{
				Message:    "Methods returning an LRO must set the metadata type.",
				Descriptor: m,
				Location:   locations.MethodOperationInfo(m),
			}}
		}

		// The netadata type should not be Empty.
		if t := lro.GetMetadataType(); t == "Empty" || t == "google.protobuf.Empty" {
			return []lint.Problem{{
				Message:    "Methods returning an LRO should create a blank message rather than using Empty.",
				Descriptor: m,
				Location:   locations.MethodOperationInfo(m),
			}}
		}

		return nil
	},
}

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
	"strings"

	"github.com/googleapis/api-linter/locations"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var lroResponse = &lint.MethodRule{
	Name:   lint.NewRuleName(151, "lro-response-type"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		lro := utils.GetOperationInfo(m)

		// Ensure the response type is set.
		if lro.GetResponseType() == "" {
			return []lint.Problem{{
				Message:    "Methods returning an LRO must set the response type.",
				Descriptor: m,
				Location:   locations.MethodOperationInfo(m),
			}}
		}

		// Unless this is a Delete method, the response type should not be Empty.
		if strings.HasPrefix(string(m.Name()), "Delete") || strings.HasPrefix(string(m.Name()), "BatchDelete") {
			return nil
		}
		if t := lro.GetResponseType(); t == "Empty" || t == "google.protobuf.Empty" {
			return []lint.Problem{{
				Message:    "Methods returning an LRO should create a blank message rather than using Empty.",
				Descriptor: m,
				Location:   locations.MethodOperationInfo(m),
			}}
		}

		return nil
	},
}

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

package aip0133

import (
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var responseLRO = &lint.MethodRule{
	Name: lint.NewRuleName(133, "response-lro"),
	OnlyIf: func(m protoreflect.MethodDescriptor) bool {
		return utils.IsCreateMethod(m) && utils.IsDeclarativeFriendlyMethod(m)
	},
	LintMethod: func(m protoreflect.MethodDescriptor) []lint.Problem {
		if !utils.IsOperation(m.Output()) {
			return []lint.Problem{{
				Message:    "Declarative-friendly create methods should use an LRO.",
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
				Suggestion: "google.longrunning.Operation",
			}}
		}
		return nil
	},
}

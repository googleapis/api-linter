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
	"reflect"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var methodSignature = &lint.MethodRule{
	Name: lint.NewRuleName(132, "method-signature"),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return isListMethod(m) && m.GetInputType().FindFieldByName("parent") != nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		signatures := utils.GetMethodSignatures(m)

		// Check if the signature is missing.
		if len(signatures) == 0 {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"List methods should include `(google.api.method_signature) = %q`",
					"parent",
				),
				Descriptor: m,
			}}
		}

		// Check if the signature is wrong.
		if !reflect.DeepEqual(signatures[0], []string{"parent"}) {
			return []lint.Problem{{
				Message:    `The method signature for List methods should be "parent".`,
				Suggestion: `option (google.api.method_signature) = "parent";`,
				Descriptor: m,
				Location:   locations.MethodSignature(m, 0),
			}}
		}
		return nil
	},
}

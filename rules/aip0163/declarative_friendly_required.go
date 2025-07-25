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

package aip0163

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var declarativeFriendlyRequired = &lint.MessageRule{
	Name: lint.NewRuleName(163, "declarative-friendly-required"),
	OnlyIf: func(m protoreflect.MessageDescriptor) bool {
		// We only want to look at request methods, not the resources themselves.
		if name := string(m.Name()); strings.HasSuffix(name, "Request") && utils.IsDeclarativeFriendlyMessage(m) {
			// If the corresponding method is a GET method, it does not need
			// validate_only.
			file, ok := m.Parent().(protoreflect.FileDescriptor)
			if !ok {
				return false
			}
			method := utils.FindMethod(file, strings.TrimSuffix(name, "Request"))
			for _, http := range utils.GetHTTPRules(method) {
				if http.Method == "GET" {
					return false
				}
			}
			return true
		}
		return false
	},
	LintMessage: func(m protoreflect.MessageDescriptor) []lint.Problem {
		if vo := m.Fields().ByName("validate_only"); vo == nil || utils.GetTypeName(vo) != "bool" || vo.IsList() {
			return []lint.Problem{{
				Message:    "Declarative-friendly mutate requests should include a singular `bool validate_only` field.",
				Descriptor: m,
			}}
		}
		return nil
	},
}

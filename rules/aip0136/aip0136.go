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

// Package aip0136 contains rules defined in https://aip.dev/136.
package aip0136

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		136,
		httpBody,
		httpMethod,
		httpNameVariable,
		httpParentVariable,
		noPrepositions,
		uriSuffix,
		verbNoun,
	)
}

func isCustomMethod(m *desc.MethodDescriptor) bool {
	// Anything with a `:` in the method URI is automatically a custom
	// method, regardless of the RPC name.
	for _, httpRule := range utils.GetHTTPRules(m) {
		if strings.Contains(httpRule.GetPlainURI(), ":") {
			return true
		}
	}

	// Methods with no `:` in the URI are standard methods if they begin with
	// one of the standard method names.
	for _, prefix := range []string{"Get", "List", "Create", "Update", "Delete", "Replace"} {
		if strings.HasPrefix(m.GetName(), prefix) {
			return false
		}
	}
	return true
}

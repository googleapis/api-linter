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
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var lroResponseReachable = &lint.MethodRule{
	Name:   lint.NewRuleName(151, "lro-response-reachable"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		return checkReachable(m, utils.GetOperationInfo(m).GetResponseType())
	},
}

func checkReachable(m *desc.MethodDescriptor, name string) []lint.Problem {
	// Ignore types defined in other packages.
	if name == "" || strings.Contains(name, ".") {
		return nil
	}

	// Make this the fully qualified type name.
	f := m.GetFile()
	if pkg := f.GetPackage(); pkg != "" {
		name = pkg + "." + name
	}

	// If the message is defined in the file, we are good to go.
	for _, file := range utils.GetAllDependencies(f) {
		if file.FindMessage(name) != nil {
			return nil
		}
	}

	// We could not find the message; complain.
	return []lint.Problem{{
		Message: fmt.Sprintf(
			"Message %q must be defined in the same file, or a file imported by this file.",
			name,
		),
		Descriptor: m,
		Location:   locations.MethodOperationInfo(m),
	}}
}

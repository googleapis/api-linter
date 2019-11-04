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
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var lroDefinedInFile = &lint.MethodRule{
	Name:   lint.NewRuleName("core", "0151", "lro-types-defined-in-file"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		lro := utils.GetOperationInfo(m)
		for k, t := range map[string]string{"Response": lro.GetResponseType(), "Metadata": lro.GetMetadataType()} {
			// Ignore types defined in other packages.
			if t == "" || strings.Contains(t, ".") {
				continue
			}

			// Complain if the message is not defined in this file.
			file := m.GetFile()
			if pkg := file.GetPackage(); pkg != "" {
				t = pkg + "." + t
			}
			if file.FindMessage(t) == nil {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"%s messages should be defined in the same file as the RPC.",
						k,
					),
					Descriptor: m,
				})
			}
		}

		return
	},
}

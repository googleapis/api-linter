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

package aip0216

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var nesting = &lint.EnumRule{
	Name: lint.NewRuleName(216, "nesting"),
	OnlyIf: func(e protoreflect.EnumDescriptor) bool {
		return strings.HasSuffix(e.Name(), "State") && e.Name() != "State"
	},
	LintEnum: func(e protoreflect.EnumDescriptor) []lint.Problem {
		messageName := strings.TrimSuffix(e.Name(), "State")
		fqMessageName := messageName
		file := e.ParentFile()
		if pkg := file.GetPackage(); pkg != "" {
			fqMessageName = pkg + "." + messageName
		}
		if file.FindMessage(fqMessageName) != nil {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Nest %q within %q, and name it `State`.",
					e.Name(),
					messageName,
				),
				Descriptor: e,
			}}
		}
		return nil
	},
}

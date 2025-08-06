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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var lroResponseReachable = &lint.MethodRule{
	Name:   lint.NewRuleName(151, "lro-response-reachable"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m protoreflect.MethodDescriptor) (problems []lint.Problem) {
		return checkReachable(m, utils.GetOperationInfo(m).GetResponseType())
	},
}

func checkReachable(m protoreflect.MethodDescriptor, name string) []lint.Problem {
	// Ignore types defined in other packages.
	if name == "" || strings.Contains(name, ".") {
		return nil
	}

	// Make this the fully qualified type name.
	f := m.ParentFile()
	if pkg := f.Package(); pkg != "" {
		name = string(pkg) + "." + name
	}

	// If the message is defined in the file, we are good to go.
	if findMessage(m.ParentFile(), protoreflect.FullName(name)) != nil {
		return nil
	}
	for _, file := range utils.GetAllDependencies(m.ParentFile()) {
		if findMessage(file, protoreflect.FullName(name)) != nil {
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

func findMessage(d protoreflect.Descriptor, name protoreflect.FullName) protoreflect.MessageDescriptor {
	switch d := d.(type) {
	case protoreflect.FileDescriptor:
		for i := 0; i < d.Messages().Len(); i++ {
			if md := findMessage(d.Messages().Get(i), name); md != nil {
				return md
			}
		}
	case protoreflect.MessageDescriptor:
		if d.FullName() == name {
			return d
		}
		for i := 0; i < d.Messages().Len(); i++ {
			if md := findMessage(d.Messages().Get(i), name); md != nil {
				return md
			}
		}
	}
	return nil
}

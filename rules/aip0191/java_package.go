// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0191

import (
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var javaPackage = &lint.FileRule{
	Name: lint.NewRuleName(191, "java-package"),
	OnlyIf: func(f protoreflect.FileDescriptor) bool {
		return hasPackage(f) && !strings.HasSuffix(f.GetPackage(), ".master")
	},
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		javaPkg := f.GetFileOptions().GetJavaPackage()
		if javaPkg == "" {
			return []lint.Problem{{
				Message:    "Proto files must set `option java_package`.",
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}
		if !strings.HasSuffix(javaPkg, f.GetPackage()) {
			return []lint.Problem{{
				Message:    "The Java Package should mirror the proto package.",
				Suggestion: fmt.Sprintf(`option java_package = "com.%s";`, f.GetPackage()),
				Descriptor: f,
				Location:   locations.FileJavaPackage(f),
			}}
		}
		return nil
	},
}

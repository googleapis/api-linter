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
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

var javaMultipleFiles = &lint.FileRule{
	Name: lint.NewRuleName(191, "java-multiple-files"),
	OnlyIf: func(f protoreflect.FileDescriptor) bool {
		return hasPackage(f) && !strings.HasSuffix(string(f.Package()), ".master")
	},
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		if !f.Options().(*dpb.FileOptions).GetJavaMultipleFiles() {
			return []lint.Problem{{
				Descriptor: f,
				Location:   locations.FilePackage(f),
				Message:    "Proto files must set `option java_multiple_files = true;`",
			}}
		}
		return nil
	},
}

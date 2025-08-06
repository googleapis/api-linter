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
	"path/filepath"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var javaOuterClassname = &lint.FileRule{
	Name: lint.NewRuleName(191, "java-outer-classname"),
	OnlyIf: func(f protoreflect.FileDescriptor) bool {
		return hasPackage(f) && !strings.HasSuffix(string(f.Package()), ".master")
	},
	LintFile: func(f protoreflect.FileDescriptor) []lint.Problem {
		filename := filepath.Base(string(f.Path()))
		want := strcase.UpperCamelCase(strings.ReplaceAll(filename, ".", "_"))

		// We ignore case on the comparisons to not be too pedantic on compound
		// word protos without underscores in the filename.
		if !strings.EqualFold(f.Options().(*descriptorpb.FileOptions).GetJavaOuterClassname(), strings.ToUpper(want)) {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Proto files should set `option java_outer_classname = %q`.",
					want,
				),
				Descriptor: f,
				Location:   locations.FilePackage(f),
			}}
		}
		return nil
	},
}

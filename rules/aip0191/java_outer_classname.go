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
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var javaOuterClassname = &lint.FileRule{
	RuleName:   lint.NewRuleName(191, "java-outer-classname"),
	OnlyIf: hasPackage,
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		if f.GetFileOptions().GetJavaOuterClassname() == "" {
			want := strcase.UpperCamelCase(strings.ReplaceAll(f.GetName(), ".", "_"))
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

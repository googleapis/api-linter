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
	"reflect"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/locations"
	"github.com/jhump/protoreflect/desc"
	"github.com/stoewer/go-strcase"
)

var consistentOptions = []string{
	"csharp_namespace",
	"go_package",
	"java_package",
	"php_class_prefix",
	"php_metadata_namespace",
	"php_namespace",
	"objc_class_prefix",
	"ruby_package",
	"swift_prefix",
}

var fileOptionConsistency = &lint.FileRule{
	Name:   lint.NewRuleName(191, "file-option-consistency"),
	OnlyIf: hasPackage,
	LintFile: func(f *desc.FileDescriptor) (problems []lint.Problem) {
		opts := f.GetFileOptions()
		for _, dep := range f.GetDependencies() {
			// We only need to look at files that are in the same package
			// as the proto we are linting.
			if dep.GetPackage() != f.GetPackage() {
				continue
			}

			// The file package options should all match between this file
			// and the file being imported.
			//
			// We will naively complain on *this* file, even though either one
			// might be the one that is wrong, and trust the API producer to do
			// the right thing.
			depOpts := dep.GetFileOptions()
			for _, opt := range consistentOptions {
				funcName := "Get" + strcase.UpperCamelCase(opt)
				a := reflect.ValueOf(opts).MethodByName(funcName).Call([]reflect.Value{})[0].String()
				b := reflect.ValueOf(depOpts).MethodByName(funcName).Call([]reflect.Value{})[0].String()
				if a != b {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Option %q should be consistent throughout the package.", opt),
						Descriptor: f,
						Location:   locations.FilePackage(f),
					})
				}
			}
		}
		return
	},
}

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
	"sort"
	"strconv"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

var consistentOptions = map[string]func(*dpb.FileOptions) string{
	"csharp_namespace":       func(o *dpb.FileOptions) string { return o.GetCsharpNamespace() },
	"go_package":             func(o *dpb.FileOptions) string { return o.GetGoPackage() },
	"java_package":           func(o *dpb.FileOptions) string { return o.GetJavaPackage() },
	"java_multiple_files":    func(o *dpb.FileOptions) string { return strconv.FormatBool(o.GetJavaMultipleFiles()) },
	"php_class_prefix":       func(o *dpb.FileOptions) string { return o.GetPhpClassPrefix() },
	"php_metadata_namespace": func(o *dpb.FileOptions) string { return o.GetPhpMetadataNamespace() },
	"php_namespace":          func(o *dpb.FileOptions) string { return o.GetPhpNamespace() },
	"objc_class_prefix":      func(o *dpb.FileOptions) string { return o.GetObjcClassPrefix() },
	"ruby_package":           func(o *dpb.FileOptions) string { return o.GetRubyPackage() },
	"swift_prefix":           func(o *dpb.FileOptions) string { return o.GetSwiftPrefix() },
}

var fileOptionConsistency = &lint.FileRule{
	Name:   lint.NewRuleName(191, "file-option-consistency"),
	OnlyIf: hasPackage,
	LintFile: func(f protoreflect.FileDescriptor) (problems []lint.Problem) {
		opts := f.Options().(*dpb.FileOptions)
		for i := 0; i < f.Imports().Len(); i++ {
			dep := f.Imports().Get(i)
			// We only need to look at files that are in the same package
			// as the proto we are linting.
			if dep.Package() != f.Package() {
				continue
			}

			// The file package options should all match between this file
			// and the file being imported.
			//
			// We will naively complain on *this* file, even though either one
			// might be the one that is wrong, and trust the API producer to do
			// the right thing.
			depOpts := dep.Options().(*dpb.FileOptions)
			for opt, valueFunc := range consistentOptions {
				if valueFunc(opts) != valueFunc(depOpts) {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Option %q should be consistent throughout the package.", opt),
						Descriptor: f,
						Location:   locations.FilePackage(f),
					})

					// Sort the problems. It does not matter for actual use, but
					// testing is hard without it since maps are iterated in randomized
					// order.
					sort.Slice(problems, func(i, j int) bool {
						return problems[i].Message < problems[j].Message
					})
				}
			}
		}
		return
	},
}

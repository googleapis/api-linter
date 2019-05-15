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

package corp

import (
	"fmt"
	"os"
	"regexp"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/descriptor"
)

var validVersionStrings = regexp.MustCompile(fmt.Sprintf(
	"%s(alpha|beta|v[\\d]+)%s",
	regexp.QuoteMeta(string(os.PathSeparator)), regexp.QuoteMeta(string(os.PathSeparator)),
))

func protoFilesMustIncludeVersion() lint.Rule {
	return &descriptor.CallbackRule{
		RuleInfo: lint.RuleInfo{
			Name:         "corp::proto_files_must_include_version",
			Description:  "All proto files must include a version number in their file path.",
			URI:          "http://go/ce-api-conformance-checks",
			RequestTypes: []lint.RequestType{lint.ProtoRequest},
		},
		Callback: descriptor.Callbacks{
			FileCallback: func(fileDescriptor protoreflect.FileDescriptor, source lint.DescriptorSource) (problems []lint.Problem, e error) {
				if !validVersionStrings.Match([]byte(fileDescriptor.Path())) {
					return []lint.Problem{{
						Message:    fmt.Sprintf("The file path %q must include a version specifier", fileDescriptor.Path()),
						Descriptor: fileDescriptor,
					}}, nil
				}

				return nil, nil
			},
		},
	}
}

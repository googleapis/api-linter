// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0215

import (
	"fmt"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

var foreignTypeReference = &lint.FieldRule{
	Name: lint.NewRuleName(202, "foreign-type-reference"),
	OnlyIf: func(fd *desc.FieldDescriptor) bool {
		return fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	},
	LintField: func(fd *desc.FieldDescriptor) []lint.Problem {
		curPkg := getPackage(fd)
		if msg := fd.GetMessageType(); msg != nil {
			msgPkg := getPackage(msg)
			if !utils.IsCommonProto(fd.GetMessageType().GetFile()) {
				// TODO: consider whether this should be less strict, and shares some common path fragment.
				// If relaxed, how much path deviation is allowed?  AIP-213 likely relates here (common components).
				if curPkg != "" && msgPkg != "" && curPkg != msgPkg {
					return []lint.Problem{{
						Message:    fmt.Sprintf("foreign type referenced, current field in %q message in %q", curPkg, msgPkg),
						Descriptor: fd,
					}}
				}
			}
		}
		return nil
	},
}

func getPackage(d desc.Descriptor) string {
	if f := d.GetFile(); f != nil {
		return f.GetPackage()
	}
	return ""
}

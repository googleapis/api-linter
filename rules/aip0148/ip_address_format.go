// Copyright 2023 Google LLC
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

package aip0148

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
)

var ipAddressFormats = []annotations.FieldInfo_Format{
	annotations.FieldInfo_IPV4,
	annotations.FieldInfo_IPV6,
	annotations.FieldInfo_IPV4_OR_IPV6,
}

var ipAddressFormat = &lint.FieldRule{
	Name: lint.NewRuleName(148, "ip-address-format"),
	OnlyIf: func(fd protoreflect.FieldDescriptor) bool {
		return fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING && (fd.Name() == "ip_address" || strings.HasSuffix(fd.Name(), "_ip_address"))
	},
	LintField: func(fd protoreflect.FieldDescriptor) []lint.Problem {
		if !utils.HasFormat(fd) || !oneofFormats(utils.GetFormat(fd), ipAddressFormats) {
			return []lint.Problem{{
				Message:    "IP Address fields must specify one of the `(google.api.field_info).format` values `IPV4`, `IPV6`, or `IPV4_OR_IPV6`",
				Descriptor: fd,
			}}
		}
		return nil
	},
}

func oneofFormats(f annotations.FieldInfo_Format, desired []annotations.FieldInfo_Format) bool {
	for _, d := range desired {
		if f == d {
			return true
		}
	}

	return false
}

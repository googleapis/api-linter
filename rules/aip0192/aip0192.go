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

// Package aip0192 contains rules defined in https://aip.dev/192.
package aip0192

import (
	"github.com/googleapis/api-linter/v2/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// AddRules adds all of the AIP-192 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		192,
		absoluteLinks,
		deprecatedComment,
		hasComments,
		noHTML,
		noMarkdownHeadings,
		noMarkdownTables,
		onlyLeadingComments,
		trademarkedNames,
	)
}

// Returns true if this is a deprecated descriptor, false otherwise.
func isDeprecated(d protoreflect.Descriptor) bool {
	switch d := d.(type) {
	case protoreflect.MethodDescriptor:
		return d.Options().(*descriptorpb.MethodOptions).GetDeprecated()
	case protoreflect.ServiceDescriptor:
		return d.Options().(*descriptorpb.ServiceOptions).GetDeprecated()
	case protoreflect.FieldDescriptor:
		return d.Options().(*descriptorpb.FieldOptions).GetDeprecated()
	case protoreflect.EnumDescriptor:
		return d.Options().(*descriptorpb.EnumOptions).GetDeprecated()
	case protoreflect.EnumValueDescriptor:
		return d.Options().(*descriptorpb.EnumValueOptions).GetDeprecated()
	case protoreflect.MessageDescriptor:
		return d.Options().(*descriptorpb.MessageOptions).GetDeprecated()
	default:
		return false
	}
}

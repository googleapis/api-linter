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
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
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

// Returns true if this is a deprecated method or service, false otherwise.
func isDeprecated(d desc.Descriptor) bool {
	switch d := d.(type) {
	case *desc.MethodDescriptor:
		return d.GetMethodOptions().GetDeprecated()
	case *desc.ServiceDescriptor:
		return d.GetServiceOptions().GetDeprecated()
	default:
		return false
	}
}

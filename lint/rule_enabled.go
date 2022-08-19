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

package lint

import (
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// defaultDisabledRules is the list of rules or groups that are by default
// disabled, because they are scoped to a very specific set of AIPs.
var defaultDisabledRules = []string{"cloud"}

// Disable all rules for deprecated descriptors.
func disableDeprecated(d desc.Descriptor) bool {
	switch v := d.(type) {
	case *desc.EnumDescriptor:
		return v.GetEnumOptions().GetDeprecated()
	case *desc.EnumValueDescriptor:
		return v.GetEnumValueOptions().GetDeprecated()
	case *desc.FieldDescriptor:
		return v.GetFieldOptions().GetDeprecated()
	case *desc.FileDescriptor:
		return v.GetFileOptions().GetDeprecated()
	case *desc.MessageDescriptor:
		return v.GetMessageOptions().GetDeprecated()
	case *desc.MethodDescriptor:
		return v.GetMethodOptions().GetDeprecated()
	case *desc.ServiceDescriptor:
		return v.GetServiceOptions().GetDeprecated()
	}
	return false
}

// Provide methods that are able to disable rules.
//
// This pattern is a hook for internal extension, by creating an additional
// file in this package that can add an additional check:
//
//	func disableForInternalReason(d desc.Descriptor) bool { ... }
//
//	func init() {
//	  descriptorDisableChecks = append(descriptorDisableChecks, disableForInternalReason)
//	}
var descriptorDisableChecks = []func(d desc.Descriptor) bool{
	disableDeprecated,
}

// ruleIsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
//
// Note, if the given source code location is not nil, it will be used to
// augment the set of commentLines.
func ruleIsEnabled(rule ProtoRule, d desc.Descriptor, l *dpb.SourceCodeInfo_Location, aliasMap map[string]string) bool {
	// If the rule is disabled because of something on the descriptor itself
	// (e.g. a deprecated annotation), address that.
	for _, mustDisable := range descriptorDisableChecks {
		// The only thing the disable functions can do is force a rule to
		// be disabled. (They can not force a rule to be enabled.)
		if mustDisable(d) {
			return false
		}
	}

	// Some rules have a legacy name. We add it to the check list.
	ruleName := string(rule.GetName())
	names := []string{ruleName, aliasMap[ruleName]}

	commentLines := []string{}
	if l != nil {
		commentLines = append(commentLines, strings.Split(l.GetLeadingComments(), "\n")...)
	}
	if f, ok := d.(*desc.FileDescriptor); ok {
		commentLines = append(commentLines, strings.Split(fileHeader(f), "\n")...)
	} else {
		commentLines = append(commentLines, strings.Split(getLeadingComments(d), "\n")...)
	}
	disabledRules := []string{}
	for _, commentLine := range commentLines {
		r := extractDisabledRuleName(commentLine)
		if r != "" {
			disabledRules = append(disabledRules, r)
		}
	}

	for _, name := range names {
		if matchRule(name, disabledRules...) {
			return false
		}
	}

	// The rule may have been disabled on a parent. (For example, a field rule
	// may be disabled at the message level to cover all fields in the message).
	//
	// Do not pass the source code location here, the source location in relation
	// to the parent is not helpful.
	if parent := d.GetParent(); parent != nil {
		return ruleIsEnabled(rule, parent, nil, aliasMap)
	}

	return true
}

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

	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// defaultDisabledRules is the list of rules or groups that are by default
// disabled, because they are scoped to a very specific set of AIPs.
var defaultDisabledRules = []string{"cloud"}

// Disable all rules for deprecated descriptors.
func disableDeprecated(d protoreflect.Descriptor) bool {
	switch v := d.(type) {
	case protoreflect.EnumDescriptor:
		return v.Options().(*dpb.EnumOptions).GetDeprecated()
	case protoreflect.EnumValueDescriptor:
		return v.Options().(*dpb.EnumValueOptions).GetDeprecated()
	case protoreflect.FieldDescriptor:
		return v.Options().(*dpb.FieldOptions).GetDeprecated()
	case protoreflect.FileDescriptor:
		return v.Options().(*dpb.FileOptions).GetDeprecated()
	case protoreflect.MessageDescriptor:
		return v.Options().(*dpb.MessageOptions).GetDeprecated()
	case protoreflect.MethodDescriptor:
		return v.Options().(*dpb.MethodOptions).GetDeprecated()
	case protoreflect.ServiceDescriptor:
		return v.Options().(*dpb.ServiceOptions).GetDeprecated()
	}
	return false
}

// Provide methods that are able to disable rules.
//
// This pattern is a hook for internal extension, by creating an additional
// file in this package that can add an additional check:
//
//	func disableForInternalReason(d protoreflect.Descriptor) bool { ... }
//
//	func init() {
//	  descriptorDisableChecks = append(descriptorDisableChecks, disableForInternalReason)
//	}
var descriptorDisableChecks = []func(d protoreflect.Descriptor) bool{
	disableDeprecated,
}

// ruleIsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
//
// Note, if the given source code location is not nil, it will be used to
// augment the set of commentLines.
func ruleIsEnabled(rule ProtoRule, d protoreflect.Descriptor, l *dpb.SourceCodeInfo_Location,
	aliasMap map[string]string, ignoreCommentDisables bool) bool {
	// If the rule is disabled because of something on the descriptor itself
	// (e.g. a deprecated annotation), address that.
	for _, mustDisable := range descriptorDisableChecks {
		// The only thing the disable functions can do is force a rule to
		// be disabled. (They can not force a rule to be enabled.)
		if mustDisable(d) {
			return false
		}
	}

	if !ignoreCommentDisables {
		if ruleIsDisabledByComments(rule, d, l, aliasMap) {
			return false
		}
	}

	// The rule may have been disabled on a parent. (For example, a field rule
	// may be disabled at the message level to cover all fields in the message).
	//
	// Do not pass the source code location here, the source location in relation
	// to the parent is not helpful.
	if parent := d.Parent(); parent != nil {
		return ruleIsEnabled(rule, parent, nil, aliasMap, ignoreCommentDisables)
	}

	return true
}

// ruleIsDisabledByComments returns true if the rule has been disabled
// by comments in the file or leading the element.
func ruleIsDisabledByComments(rule ProtoRule, d protoreflect.Descriptor, l *dpb.SourceCodeInfo_Location, aliasMap map[string]string) bool {
	// Some rules have a legacy name. We add it to the check list.
	ruleName := string(rule.GetName())
	names := []string{ruleName, aliasMap[ruleName]}

	commentLines := []string{}
	if l != nil {
		commentLines = append(commentLines, strings.Split(l.GetLeadingComments(), "\n")...)
	}
	if f, ok := d.(protoreflect.FileDescriptor); ok {
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
			return true
		}
	}

	return false
}

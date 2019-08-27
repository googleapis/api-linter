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
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule struct {
	// Name contains the rule's name.
	Name RuleName

	// Description is a short description of the rule,
	// used in documentation and elsewhere.
	Description string

	// URI is the address where the guideline is documented.
	// This should be displayed to give API designers more information about
	// "how to do this right".
	URI string

	// LintFile is called for files as a whole.
	LintFile func(*desc.FileDescriptor) []Problem

	// LintMessage is called for any messages in the file.
	// This includes nested messages, regardless of depth.
	LintMessage func(*desc.MessageDescriptor) []Problem

	// LintField is called for any fields in the file.
	// This includes fields within nested messages, regardless of depth.
	LintField func(*desc.FieldDescriptor) []Problem

	// LintService is called for each service in the file.
	LintService func(*desc.ServiceDescriptor) []Problem

	// LintMethod is called for each method in the file.
	LintMethod func(*desc.MethodDescriptor) []Problem

	// LintEnum is called for each enum in the file.
	LintEnum func(*desc.EnumDescriptor) []Problem
}

// Lint iterates over every message in a FileDescriptor and runs the
// callback, aggregating any problems that are found.
func (rule *Rule) Lint(f *desc.FileDescriptor) (problems []Problem) {
	// Iterate over each message and process for any kind of rules within
	// the message.
	//
	// Note: Messages can contain enums in addition to messages and fields, so
	// enums are processed both here *and* at the top level.
	for _, message := range getAllMessages(f) {
		if rule.LintMessage != nil {
			problems = append(problems, rule.LintMessage(message)...)
		}
		for _, field := range message.GetFields() {
			if rule.LintField != nil {
				problems = append(problems, rule.LintField(field)...)
			}
		}
		for _, enum := range message.GetNestedEnumTypes() {
			if rule.LintEnum != nil {
				problems = append(problems, rule.LintEnum(enum)...)
			}
		}
	}

	// Iterate over each service and process any rules within that service,
	// as well as rules for each method.
	for _, service := range f.GetServices() {
		if rule.LintService != nil {
			problems = append(problems, rule.LintService(service)...)
		}
		for _, method := range service.GetMethods() {
			if rule.LintMethod != nil {
				problems = append(problems, rule.LintMethod(method)...)
			}
		}
	}

	// Process rules for top-level enums.
	for _, enum := range f.GetEnumTypes() {
		if rule.LintEnum != nil {
			problems = append(problems, rule.LintEnum(enum)...)
		}
	}

	// Finally, process rules for the file itself.
	if rule.LintFile != nil {
		problems = append(problems, rule.LintFile(f)...)
	}

	// Done aggregating problems for this file; return the list.
	return problems
}

// IsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
func (rule *Rule) IsEnabled(d desc.Descriptor) bool {
	directive := fmt.Sprintf("api-linter: %s=disabled", rule.Name)

	// If the comments above the descriptor disable the rule,
	// return true.
	if sourceInfo := d.GetSourceInfo(); sourceInfo != nil {
		if strings.Contains(sourceInfo.GetLeadingComments(), directive) {
			return false
		}
	}

	// The rule may also be disabled at the file level; if it is, return true.
	if sourceInfo := d.GetFile().GetSourceInfo(); sourceInfo != nil {
		for _, line := range sourceInfo.GetLeadingDetachedComments() {
			if strings.Contains(line, directive) {
				return false
			}
		}
	}

	// The rule is enabled.
	return true
}

// getAllMessages returns a slice with every message (not just top-level
// messages) in the file.
func getAllMessages(f *desc.FileDescriptor) (messages []*desc.MessageDescriptor) {
	for _, message := range f.GetMessageTypes() {
		messages = append(messages, message)
		messages = append(messages, getAllNestedMessages(message)...)
	}
	return messages
}

// getAllNestedMessages returns a slice with the given message descriptor as well
// as all nested message descriptors, traversing to arbitrary depth.
func getAllNestedMessages(m *desc.MessageDescriptor) (messages []*desc.MessageDescriptor) {
	for _, nested := range m.GetNestedMessageTypes() {
		messages = append(messages, nested)
		messages = append(messages, getAllNestedMessages(nested)...)
	}
	return messages
}

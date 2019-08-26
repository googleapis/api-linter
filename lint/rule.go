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

import "github.com/jhump/protoreflect/desc"

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule struct {
	// Info contains metadata about a rule.
	Info RuleInfo

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

	// Finally, process rules for top-level enums.
	for _, enum := range f.GetEnumTypes() {
		if rule.LintEnum != nil {
			problems = append(problems, rule.LintEnum(enum)...)
		}
	}

	// Done aggregating problems for this file; return the list.
	return problems
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

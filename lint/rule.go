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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// Rule defines a lint rule that checks Google Protobuf APIs.
type Rule interface {
	// GetName returns the name of the rule.
	GetName() RuleName

	// GetURI returns the URI where the applicable guideline
	// is documented. (This should generally be an AIP on https://aip.dev/.)
	GetURI() string

	// Lint accepts a FileDescriptor and lints it,
	// returning a slice of Problem objects it finds.
	Lint(*desc.FileDescriptor) []Problem
}

// FileRule defines a lint rule that checks a file as a whole.
type FileRule struct {
	Name RuleName
	URI  string

	// LintFile accepts a FileDescriptor and lints it, returning a slice of
	// Problems it finds.
	LintFile func(*desc.FileDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *FileRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *FileRule) GetURI() string {
	return r.URI
}

// Lint forwards the FileDescriptor to the LintFile method defined on the
// FileRule.
func (r *FileRule) Lint(fd *desc.FileDescriptor) []Problem {
	return r.LintFile(fd)
}

// MessageRule defines a lint rule that is run on each message (top-level or
// nested) within a file.
type MessageRule struct {
	Name RuleName
	URI  string

	// LintMessage accepts a MessageDescriptor and lints it, returning a slice
	// of Problems it finds.
	LintMessage func(*desc.MessageDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *MessageRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *MessageRule) GetURI() string {
	return r.URI
}

// Lint accepts a FileDescriptor and iterates over every message in the
// file, and lints each message in the file.
func (r *MessageRule) Lint(fd *desc.FileDescriptor) []Problem {
	problems := []Problem{}

	// Iterate over each message and process rules for each message.
	for _, message := range getAllMessages(fd) {
		problems = append(problems, r.LintMessage(message)...)
	}
	return problems
}

// FieldRule defines a lint rule that is run on each field within a file.
type FieldRule struct {
	Name RuleName
	URI  string

	// LintField accepts a FieldDescriptor and lints it, returning a slice of
	// Problems it finds.
	LintField func(*desc.FieldDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *FieldRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *FieldRule) GetURI() string {
	return r.URI
}

// Lint accepts a FileDescriptor and lints every field in the file.
func (r *FieldRule) Lint(fd *desc.FileDescriptor) (problems []Problem) {
	// Iterate over each message and process rules for each field in that
	// message.
	for _, message := range getAllMessages(fd) {
		for _, field := range message.GetFields() {
			problems = append(problems, r.LintField(field)...)
		}
	}
	return problems
}

// ServiceRule defines a lint rule that is run on each service.
type ServiceRule struct {
	Name RuleName
	URI  string

	// LintService accepts a ServiceDescriptor and lints it.
	LintService func(*desc.ServiceDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *ServiceRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *ServiceRule) GetURI() string {
	return r.URI
}

// Lint accepts a FileDescriptor and lints every service in the file.
func (r *ServiceRule) Lint(fd *desc.FileDescriptor) (problems []Problem) {
	for _, service := range fd.GetServices() {
		problems = append(problems, r.LintService(service)...)
	}
	return problems
}

// MethodRule defines a lint rule that is run on each method.
type MethodRule struct {
	Name RuleName
	URI  string

	// LintMethod accepts a MethodDescriptor and lints it.
	LintMethod func(*desc.MethodDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *MethodRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *MethodRule) GetURI() string {
	return r.URI
}

// Lint accepts a FileDescriptor and lints every method in the file.
func (r *MethodRule) Lint(fd *desc.FileDescriptor) (problems []Problem) {
	for _, service := range fd.GetServices() {
		for _, method := range service.GetMethods() {
			problems = append(problems, r.LintMethod(method)...)
		}
	}
	return problems
}

// EnumRule defines a lint rule that is run on each enum.
type EnumRule struct {
	Name RuleName
	URI  string

	// LintEnum accepts a EnumDescriptor and lints it.
	LintEnum func(*desc.EnumDescriptor) []Problem

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *EnumRule) GetName() RuleName {
	return r.Name
}

// GetURI returns the URI where the applicable guideline is documented.
func (r *EnumRule) GetURI() string {
	return r.URI
}

// Lint accepts a FileDescriptor and lints every enum in the file
// (including enums nested within messages).
func (r *EnumRule) Lint(fd *desc.FileDescriptor) (problems []Problem) {
	// Lint enums that are at the top level of the file.
	for _, enum := range fd.GetEnumTypes() {
		problems = append(problems, r.LintEnum(enum)...)
	}

	// Lint enums that are nested within messages.
	for _, message := range getAllMessages(fd) {
		for _, enum := range message.GetNestedEnumTypes() {
			problems = append(problems, r.LintEnum(enum)...)
		}
	}
	return problems
}

// ruleIsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
func ruleIsEnabled(rule Rule, d desc.Descriptor) bool {
	directive := fmt.Sprintf("api-linter: %s=disabled", rule.GetName())

	// If the comments above the descriptor disable the rule,
	// return false.
	if sourceInfo := d.GetSourceInfo(); sourceInfo != nil {
		if strings.Contains(sourceInfo.GetLeadingComments(), directive) {
			return false
		}
	}

	// The rule may also be disabled at the file level.
	// If it is, return false.
	if strings.Contains(fileHeader(d.GetFile()), directive) {
		return false
	}

	// The rule is enabled.
	return true
}

// getAllMessages returns a slice with every message (not just top-level
// messages) in the file.
func getAllMessages(f *desc.FileDescriptor) (messages []*desc.MessageDescriptor) {
	messages = append(messages, f.GetMessageTypes()...)
	for _, message := range f.GetMessageTypes() {
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

// fileHeader attempts to get the comment at the top of the file, but it
// is on a best effort basis because protobuf is inconsistent.
//
// Taken from https://github.com/jhump/protoreflect/issues/215
func fileHeader(fd *desc.FileDescriptor) string {
	var firstLoc *descriptor.SourceCodeInfo_Location
	var firstSpan int64
	for _, curr := range fd.AsFileDescriptorProto().GetSourceCodeInfo().GetLocation() {
		if curr.LeadingComments == nil && len(curr.LeadingDetachedComments) == 0 {
			// Skip locations that have no comments.
			continue
		}
		currSpan := asPos(curr.Span)
		if firstLoc == nil || currSpan < firstSpan {
			firstLoc = curr
			firstSpan = currSpan
		}
	}
	if firstLoc == nil {
		return ""
	}
	if len(firstLoc.LeadingDetachedComments) > 0 {
		return strings.Join(firstLoc.LeadingDetachedComments, "\n")
	}
	return firstLoc.GetLeadingComments()
}

func asPos(span []int32) int64 {
	return (int64(span[0]) << 32) + int64(span[1])
}

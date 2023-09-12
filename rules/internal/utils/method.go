// Copyright 2023 Google LLC
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

package utils

import (
	"regexp"

	"github.com/jhump/protoreflect/desc"
)

var (
	createMethodRegexp         = regexp.MustCompile("^Create(?:[A-Z]|$)")
	getMethodRegexp            = regexp.MustCompile("^Get(?:[A-Z]|$)")
	listMethodRegexp           = regexp.MustCompile("^List(?:[A-Z]|$)")
	listRevisionsMethodRegexp  = regexp.MustCompile(`^List(?:[A-Za-z0-9]+)Revisions$`)
	updateMethodRegexp         = regexp.MustCompile("^Update(?:[A-Z]|$)")
	deleteMethodRegexp         = regexp.MustCompile("^Delete(?:[A-Z]|$)")
	deleteRevisionMethodRegexp = regexp.MustCompile("^Delete[A-Za-z0-9]*Revision$")
)

// IsCreateMethod returns true if this is a AIP-133 Create method.
func IsCreateMethod(m *desc.MethodDescriptor) bool {
	return createMethodRegexp.MatchString(m.GetName())
}

// IsGetMethod returns true if this is a AIP-131 Get method.
func IsGetMethod(m *desc.MethodDescriptor) bool {
	methodName := m.GetName()
	if methodName == "GetIamPolicy" {
		return false
	}
	return getMethodRegexp.MatchString(methodName)
}

// IsListMethod return true if this is an AIP-132 List method
func IsListMethod(m *desc.MethodDescriptor) bool {
	return listMethodRegexp.MatchString(m.GetName()) && !IsListRevisionsMethod(m)
}

// IsListRevisionsMethod returns true if this is an AIP-162 List Revisions method,
// false otherwise.
func IsListRevisionsMethod(m *desc.MethodDescriptor) bool {
	return listRevisionsMethodRegexp.MatchString(m.GetName())
}

// IsUpdateMethod returns true if this is a AIP-134 Update method
func IsUpdateMethod(m *desc.MethodDescriptor) bool {
	methodName := m.GetName()
	return updateMethodRegexp.MatchString(methodName)
}

// Returns true if this is a AIP-135 Delete method, false otherwise.
func IsDeleteMethod(m *desc.MethodDescriptor) bool {
	return deleteMethodRegexp.MatchString(m.GetName()) && !deleteRevisionMethodRegexp.MatchString(m.GetName())
}

// GetListResourceMessage returns the resource for a list method,
// nil otherwise.
func GetListResourceMessage(m *desc.MethodDescriptor) *desc.MessageDescriptor {
	repeated := GetRepeatedMessageFields(m.GetOutputType())
	if len(repeated) > 0 {
		return repeated[0].GetMessageType()
	}
	return nil
}

// IsStreaming returns if the method is either client or server streaming.
func IsStreaming(m *desc.MethodDescriptor) bool {
	return m.IsClientStreaming() || m.IsServerStreaming()
}

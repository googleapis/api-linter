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

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	createMethodRegexp   = regexp.MustCompile("^Create(?:[A-Z]|$)")
	getMethodRegexp      = regexp.MustCompile("^Get(?:[A-Z]|$)")
	listMethodRegexp     = regexp.MustCompile("^List(?:[A-Z]|$)")
	updateMethodRegexp   = regexp.MustCompile("^Update(?:[A-Z]|$)")
	deleteMethodRegexp   = regexp.MustCompile("^Delete(?:[A-Z]|$)")
	standardMethodRegexp = regexp.MustCompile("^(Batch(Get|Create|Update|Delete))|(Get|Create|Update|Delete|List)(?:[A-Z]|$)")

	// AIP-162 Resource revision methods
	listRevisionsMethodRegexp        = regexp.MustCompile(`^List(?:[A-Za-z0-9]+)Revisions$`)
	legacyListRevisionsURINameRegexp = regexp.MustCompile(`:listRevisions$`)
	commitRevisionMethodRegexp       = regexp.MustCompile(`^Commit([A-Za-z0-9]+)$`)
	deleteRevisionMethodRegexp       = regexp.MustCompile(`^Delete([A-Za-z0-9]+)Revision$`)
	rollbackRevisionMethodRegexp     = regexp.MustCompile(`^Rollback([A-Za-z0-9]+)$`)
	tagRevisionMethodRegexp          = regexp.MustCompile(`^Tag([A-Za-z0-9]+)Revision$`)
)

// IsCreateMethod returns true if this is a AIP-133 Create method.
func IsCreateMethod(m protoreflect.MethodDescriptor) bool {
	return createMethodRegexp.MatchString(string(m.Name()))
}

// IsCreateMethodWithResolvedReturnType returns true if this is a AIP-133 Create method with
// a non-nil response type. This method should be used for filtering in linter
// rules which access the response type of the method, to avoid crashing due
// to dereferencing a nil pointer to the response.
func IsCreateMethodWithResolvedReturnType(m protoreflect.MethodDescriptor) bool {
	if !IsCreateMethod(m) {
		return false
	}

	return GetResponseType(m) != nil
}

// IsGetMethod returns true if this is a AIP-131 Get method.
func IsGetMethod(m protoreflect.MethodDescriptor) bool {
	methodName := string(m.Name())
	if methodName == "GetIamPolicy" {
		return false
	}
	return getMethodRegexp.MatchString(methodName)
}

// IsListMethod return true if this is an AIP-132 List method
func IsListMethod(m protoreflect.MethodDescriptor) bool {
	return listMethodRegexp.MatchString(string(m.Name())) && !IsLegacyListRevisionsMethod(m)
}

// IsLegacyListRevisions identifies such a method by having the appropriate
// method name, having a `name` field instead of parent, and a HTTP suffix of
// `listRevisions`.
func IsLegacyListRevisionsMethod(m protoreflect.MethodDescriptor) bool {
	// Must be named like List{Resource}Revisions.
	if !listRevisionsMethodRegexp.MatchString(string(m.Name())) {
		return false
	}

	// Must have a `name` field instead of `parent`.
	if m.Input().Fields().ByName("name") == nil {
		return false
	}

	// Must have the `:listRevisions` HTTP URI suffix.
	if !HasHTTPRules(m) {
		// If it doesn't have HTTP bindings, we shouldn't proceed to the next
		// check, but a List{Resource}Revisions method with a `name` field is
		// probably enough to be sure in the absence of HTTP bindings.
		return true
	}

	// Just check the first bidning as they should all have the same suffix.
	h := GetHTTPRules(m)[0].GetPlainURI()
	return legacyListRevisionsURINameRegexp.MatchString(h)
}

// IsUpdateMethod returns true if this is a AIP-134 Update method
func IsUpdateMethod(m protoreflect.MethodDescriptor) bool {
	methodName := string(m.Name())
	return updateMethodRegexp.MatchString(methodName)
}

// Returns true if this is a AIP-135 Delete method, false otherwise.
func IsDeleteMethod(m protoreflect.MethodDescriptor) bool {
	return deleteMethodRegexp.MatchString(string(m.Name())) && !deleteRevisionMethodRegexp.MatchString(string(m.Name()))
}

// GetListResourceMessage returns the resource for a list method,
// nil otherwise.
func GetListResourceMessage(m protoreflect.MethodDescriptor) protoreflect.MessageDescriptor {
	repeated := GetRepeatedMessageFields(m.Output())
	if len(repeated) > 0 {
		return repeated[0].Message()
	}
	return nil
}

// IsStreaming returns if the method is either client or server streaming.
func IsStreaming(m protoreflect.MethodDescriptor) bool {
	return m.IsStreamingClient() || m.IsStreamingServer()
}

// IsStandardMethod returns true if this is a AIP-130 Standard Method
func IsStandardMethod(m protoreflect.MethodDescriptor) bool {
	return standardMethodRegexp.MatchString(string(m.Name()))
}

// IsCustomMethod returns true if this is a AIP-136 Custom Method
func IsCustomMethod(m protoreflect.MethodDescriptor) bool {
	return !IsStandardMethod(m) && !isRevisionMethod(m)
}

// isRevisionMethod returns true if the given method is any of the documented
// Revision methods. At the moment, this is only relevant for excluding revision
// methods from other method type checks.
func isRevisionMethod(m protoreflect.MethodDescriptor) bool {
	return IsDeleteRevisionMethod(m) ||
		IsTagRevisionMethod(m) ||
		IsCommitRevisionMethod(m) ||
		IsRollbackRevisionMethod(m)
}

// IsDeleteRevisionMethod returns true if this is an AIP-162 Delete Revision
// method, false otherwise.
func IsDeleteRevisionMethod(m protoreflect.MethodDescriptor) bool {
	return deleteRevisionMethodRegexp.MatchString(string(m.Name()))
}

// IsTagRevisionMethod returns true if this is an AIP-162 Tag Revision method,
// false otherwise.
func IsTagRevisionMethod(m protoreflect.MethodDescriptor) bool {
	return tagRevisionMethodRegexp.MatchString(string(m.Name()))
}

// IsCommitRevisionMethod returns true if this is an AIP-162 Commit method,
// false otherwise.
func IsCommitRevisionMethod(m protoreflect.MethodDescriptor) bool {
	return commitRevisionMethodRegexp.MatchString(string(m.Name()))
}

// IsRollbackRevisionMethod returns true if this is an AIP-162 Rollback method,
// false otherwise.
func IsRollbackRevisionMethod(m protoreflect.MethodDescriptor) bool {
	return rollbackRevisionMethodRegexp.MatchString(string(m.Name()))
}

// ExtractRevisionResource uses the appropriate revision method regular
// expression to capture and extract the resource noun in the method name.
// If the given method is not one of the standard revision RPCs, it returns
// empty string and false.
func ExtractRevisionResource(m protoreflect.MethodDescriptor) (string, bool) {
	if !isRevisionMethod(m) {
		return "", false
	}

	n := string(m.Name())

	if IsCommitRevisionMethod(m) {
		return commitRevisionMethodRegexp.FindStringSubmatch(n)[1], true
	} else if IsTagRevisionMethod(m) {
		return tagRevisionMethodRegexp.FindStringSubmatch(n)[1], true
	} else if IsRollbackRevisionMethod(m) {
		return rollbackRevisionMethodRegexp.FindStringSubmatch(n)[1], true
	} else if IsDeleteRevisionMethod(m) {
		return deleteRevisionMethodRegexp.FindStringSubmatch(n)[1], true
	}

	return "", false
}
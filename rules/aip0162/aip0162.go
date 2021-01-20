// Copyright 2021 Google LLC
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

// Package aip0162 contains rules defined in https://aip.dev/162.
package aip0162

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		162,
		commitHTTPURISuffix,
		tagRevisionHTTPBody,
		tagRevisionHTTPMethod,
		tagRevisionHTTPURISuffix,
		tagRevisionRequestMessageName,
		tagRevisionRequestNameBehavior,
		tagRevisionRequestNameField,
		tagRevisionRequestNameReference,
		tagRevisionRequestTagBehavior,
		tagRevisionRequestTagField,
		tagRevisionResponseMessageName,
	)
}

var tagRevisionMethodRegexp = regexp.MustCompile(`^Tag([A-Za-z0-9]+)Revision$`)
var tagRevisionReqMessageRegexp = regexp.MustCompile(`^Tag(?:[A-Za-z0-9]+)RevisionRequest$`)
var tagRevisionURINameRegexp = regexp.MustCompile(`:tagRevision$`)

// Returns true if this is an AIP-162 Tag Revision method, false otherwise.
func isTagRevisionMethod(m *desc.MethodDescriptor) bool {
	return tagRevisionMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AIP-162 Tag Revision request message, false otherwise.
func isTagRevisionRequestMessage(m *desc.MessageDescriptor) bool {
	return tagRevisionReqMessageRegexp.MatchString(m.GetName())
}

var commitMethodRegexp = regexp.MustCompile(`^Commit(?:[A-Za-z0-9]+)$`)
var commitURINameRegexp = regexp.MustCompile(`:commit$`)

// Returns true if this is an AIP-162 Commit method, false otherwise.
func isCommitMethod(m *desc.MethodDescriptor) bool {
	return commitMethodRegexp.MatchString(m.GetName())
}

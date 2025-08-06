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

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		162,
		commitHTTPBody,
		commitHTTPMethod,
		commitHTTPURISuffix,
		commitRequestMessageName,
		commitRequestNameBehavior,
		commitRequestNameField,
		commitRequestNameReference,
		commitResponseMessageName,
		deleteRevisionHTTPBody,
		deleteRevisionHTTPMethod,
		deleteRevisionHTTPURISuffix,
		deleteRevisionRequestMessageName,
		deleteRevisionRequestNameBehavior,
		deleteRevisionRequestNameField,
		deleteRevisionRequestNameReference,
		deleteRevisionResponseMessageName,
		rollbackHTTPBody,
		rollbackHTTPMethod,
		rollbackHTTPURISuffix,
		rollbackRequestMessageName,
		rollbackRequestNameBehavior,
		rollbackRequestNameField,
		rollbackRequestNameReference,
		rollbackRequestRevisionIDBehavior,
		rollbackRequestRevisionIDField,
		rollbackResponseMessageName,
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

var (
	tagRevisionReqMessageRegexp = regexp.MustCompile(`^Tag(?:[A-Za-z0-9]+)RevisionRequest$`)
	tagRevisionURINameRegexp    = regexp.MustCompile(`:tagRevision$`)
)

// Returns true if this is an AIP-162 Tag Revision request message, false otherwise.
func isTagRevisionRequestMessage(m protoreflect.MessageDescriptor) bool {
	return tagRevisionReqMessageRegexp.MatchString(string(m.Name()))
}

var (
	commitReqMessageRegexp = regexp.MustCompile(`^Commit(?:[A-Za-z0-9]+)Request$`)
	commitURINameRegexp    = regexp.MustCompile(`:commit$`)
)

// Returns true if this is an AIP-162 Commit request message, false otherwise.
func isCommitRequestMessage(m protoreflect.MessageDescriptor) bool {
	return commitReqMessageRegexp.MatchString(string(m.Name()))
}

var (
	rollbackReqMessageRegexp = regexp.MustCompile(`^Rollback(?:[A-Za-z0-9]+)Request$`)
	rollbackURINameRegexp    = regexp.MustCompile(`:rollback$`)
)

// Returns true if this is an AIP-162 Rollback request message, false otherwise.
func isRollbackRequestMessage(m protoreflect.MessageDescriptor) bool {
	return rollbackReqMessageRegexp.MatchString(string(m.Name()))
}

var (
	deleteRevisionReqMessageRegexp = regexp.MustCompile(`^Delete(?:[A-Za-z0-9]+)RevisionRequest$`)
	deleteRevisionURINameRegexp    = regexp.MustCompile(`:deleteRevision$`)
)

// Returns true if this is an AIP-162 Delete Revision request message, false otherwise.
func isDeleteRevisionRequestMessage(m protoreflect.MessageDescriptor) bool {
	return deleteRevisionReqMessageRegexp.MatchString(string(m.Name()))
}

// IsListRevisionsRequestMessage returns true if this is an AIP-162 List
// Revisions request message, false otherwise.
//
// Deprecated: Use the same method from the utils package instead.
func IsListRevisionsRequestMessage(m protoreflect.MessageDescriptor) bool {
	return utils.IsListRevisionsRequestMessage(m)
}

// IsListRevisionsResponseMessage returns true if this is an AIP-162 List
// Revisions response message, false otherwise.
//
// Deprecated: Use the same method from the utils package instead.
func IsListRevisionsResponseMessage(m protoreflect.MessageDescriptor) bool {
	return utils.IsListRevisionsResponseMessage(m)
}
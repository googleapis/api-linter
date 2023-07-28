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
	listReqMessageRegexp           = regexp.MustCompile("^List[A-Za-z0-9]*Request$")
	listRespMessageRegexp          = regexp.MustCompile("^List([A-Za-z0-9]*)Response$")
	listRevisionsReqMessageRegexp  = regexp.MustCompile(`^List(?:[A-Za-z0-9]+)RevisionsRequest$`)
	listRevisionsRespMessageRegexp = regexp.MustCompile(`^List(?:[A-Za-z0-9]+)RevisionsResponse$`)
)

// Return true if this is an AIP-132 List request message, false otherwise.
func IsListRequestMessage(m *desc.MessageDescriptor) bool {
	return listReqMessageRegexp.MatchString(m.GetName()) && !IsListRevisionsRequestMessage(m)
}

// Return true if this is an AIP-132 List response message, false otherwise.
func IsListResponseMessage(m *desc.MessageDescriptor) bool {
	return listRespMessageRegexp.MatchString(m.GetName()) && !IsListRevisionsResponseMessage(m)
}

// IsListRevisionsRequestMessage returns true if this is an AIP-162 List
// Revisions request message, false otherwise.
func IsListRevisionsRequestMessage(m *desc.MessageDescriptor) bool {
	return listRevisionsReqMessageRegexp.MatchString(m.GetName())
}

// IsListRevisionsResponseMessage returns true if this is an AIP-162 List
// Revisions response message, false otherwise.
func IsListRevisionsResponseMessage(m *desc.MessageDescriptor) bool {
	return listRevisionsRespMessageRegexp.MatchString(m.GetName())
}

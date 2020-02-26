// Copyright 2019 Google LLC
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

// Package aip0158 contains rules defined in https://aip.dev/158.
package aip0158

import (
	"regexp"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AIP-158 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		158,
		requestPaginationPageSize,
		requestPaginationPageToken,
		responsePaginationNextPageToken,
		responsePluralFirstField,
	)
}

var paginatedReq = regexp.MustCompile("^(List|Search)[A-Za-z0-9]*Request$")
var paginatedRes = regexp.MustCompile("^(List|Search)[A-Za-z0-9]*Response$")

// Return true if this is an AIP-158 List request message, false otherwise.
func isPaginatedRequestMessage(m *desc.MessageDescriptor) bool {
	return paginatedReq.MatchString(m.GetName()) || m.FindFieldByName("page_size") != nil || m.FindFieldByName("page_token") != nil
}

// Return true if this is an AIP-158 List response message, false otherwise.
func isPaginatedResponseMessage(m *desc.MessageDescriptor) bool {
	return paginatedRes.MatchString(m.GetName()) || m.FindFieldByName("next_page_token") != nil
}

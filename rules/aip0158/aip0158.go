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
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// AddRules adds all of the AIP-158 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		158,
		requestPaginationPageSize,
		requestPaginationPageToken,
		requestSkipField,
		responsePaginationNextPageToken,
		responseRepeatedFirstField,
		responsePluralFirstField,
		responseUnary,
	)
}

var (
	paginatedReq = regexp.MustCompile("^(List|Search)[A-Za-z0-9]*Request$")
	paginatedRes = regexp.MustCompile("^(List|Search)[A-Za-z0-9]*Response$")
)

// Return true if this is an AIP-158 List request message, false otherwise.
func isPaginatedRequestMessage(m protoreflect.MessageDescriptor) bool {
	if paginatedReq.MatchString(string(m.Name())) {
		return true
	}
	// Ignore messages that happen to have these fields but are not requests.
	if !strings.HasSuffix(string(m.Name()), "Request") {
		return false
	}
	return m.Fields().ByName("page_size") != nil || m.Fields().ByName("page_token") != nil
}

// Return true if this is an AIP-158 List response message, false otherwise.
func isPaginatedResponseMessage(m protoreflect.MessageDescriptor) bool {
	return paginatedRes.MatchString(string(m.Name())) || m.Fields().ByName("next_page_token") != nil
}

func isPaginatedMethod(m protoreflect.MethodDescriptor) bool {
	return isPaginatedRequestMessage(m.Input()) && isPaginatedResponseMessage(m.Output())
}

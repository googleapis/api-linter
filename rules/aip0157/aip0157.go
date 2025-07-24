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

// Package aip0157 contains rules defined in https://aip.dev/157.
package aip0157

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// AddRules accepts a register function and registers each of this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		157,
		requestReadMaskField,
	)
}

func isRequestMessage(m protoreflect.MessageDescriptor) bool {
	return strings.HasSuffix(m.Name(), "Request")
}

// Copyright 2020 Google LLC
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

// Package aip0128 contains rules defined in https://aip.dev/128.
package aip0128

import (
	"github.com/commure/api-linter/lint"
	"github.com/commure/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		128,
		resourceAnnotationsField,
		resourceReconcilingBehavior,
		resourceReconcilingField,
	)
}

func isDeclarativeFriendlyResource(m *desc.MessageDescriptor) bool {
	// IsDeclarativeFriendly returns true for both
	// resources and request messages, but we only care about resources.
	resource := utils.DeclarativeFriendlyResource(m)
	return resource != nil && resource == m
}

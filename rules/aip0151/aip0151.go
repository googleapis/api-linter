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

// Package aip0151 contains rules defined in https://aip.dev/151.
package aip0151

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AIP-151 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		151,
		lroAnnotationExists,
		lroDefinedInFile,
		lroMetadata,
		lroResponse,
	)
}

func isLRO(m *desc.MethodDescriptor) bool {
	return m.GetOutputType().GetFullyQualifiedName() == "google.longrunning.Operation"
}

func isAnnotatedLRO(m *desc.MethodDescriptor) bool {
	return isLRO(m) && utils.GetOperationInfo(m) != nil
}

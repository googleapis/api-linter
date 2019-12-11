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

package aip0151

import (
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var lroMetadataReachable = &lint.MethodRule{
	Name:   lint.NewRuleName(151, "lro-metadata-reachable"),
	OnlyIf: isAnnotatedLRO,
	LintMethod: func(m *desc.MethodDescriptor) (problems []lint.Problem) {
		// See lro_response_reachable.go for `checkReachable` method.
		return checkReachable(m, utils.GetOperationInfo(m).GetMetadataType())
	},
}

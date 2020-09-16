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

// Package aip0124 contains rules defined in https://aip.dev/124.
package aip0124

import (
	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AIP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		124,
		referenceSamePackage,
		validReference,
	)
}

// isUnknownType returns true if and only if the type is not a common,
// well known type.
func isUnknownType(f *desc.FieldDescriptor) bool {
	if ref := utils.GetResourceReference(f); ref != nil {
		urt := ref.GetType()
		if urt == "" {
			urt = ref.GetChildType()
		}
		return !stringset.New(
			// Allow the common resource types in GCP.
			// FIXME: Modularize this.
			"cloudresourcemanager.googleapis.com/Project",
			"cloudresourcemanager.googleapis.com/Organization",
			"cloudresourcemanager.googleapis.com/Folder",
			"billing.googleapis.com/BillingAccount",
			"locations.googleapis.com/Location",

			// Allow *.
			"*",

			// If no type is declared, ignore this.
			"",
		).Contains(urt)
	}
	return false
}

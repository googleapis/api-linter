// Copyright 2023 Google LLC
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

package utils

import (
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

// GetResourceSingular returns the resource singular. The
// empty string is returned if the singular cannot be found.
// Since the singular is not always annotated, it extracts
// it from multiple different locations including:
// 1. the singular annotation
// 2. the type definition
func GetResourceSingular(r *apb.ResourceDescriptor) string {
	if r == nil {
		return ""
	}
	if r.Singular != "" {
		return r.Singular
	}
	if r.Type != "" {
		_, tn, ok := SplitResourceTypeName(r.Type)
		if ok {
			return tn
		}
	}
	return ""
}

// GetResourcePlural is a convenience method for getting the `plural` field of a
// resource.
func GetResourcePlural(r *apb.ResourceDescriptor) string {
	if r == nil {
		return ""
	}

	return r.Plural
}

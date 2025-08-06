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
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetResourceSingular returns the resource singular. The
// empty string is returned if the singular cannot be found.
// Since the singular is not always annotated, it extracts
// it from multiple different locations including:
// 1. the singular annotation
// 2. the type definition
func GetResourceSingular(r *annotations.ResourceDescriptor) string {
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
func GetResourcePlural(r *annotations.ResourceDescriptor) string {
	if r == nil {
		return ""
	}

	return r.Plural
}

// GetResourceNameField is a convenience method for getting the name of the
// field that represents the resource's name. This is either set by the
// `name_field` attribute, or defaults to "name".
func GetResourceNameField(r *annotations.ResourceDescriptor) string {
	if r == nil {
		return ""
	}
	if n := r.GetNameField(); n != "" {
		return n
	}
	return "name"
}

// IsResourceRevision determines if the given message represents a resource
// revision as described in AIP-162.
func IsResourceRevision(m protoreflect.MessageDescriptor) bool {
	return IsResource(m) && strings.HasSuffix(string(m.Name()), "Revision")
}

// IsRevisionRelationship determines if the "revision" resource is actually
// a revision of the "parent" resource.
func IsRevisionRelationship(parent, revision *annotations.ResourceDescriptor) bool {
	_, pType, ok := SplitResourceTypeName(parent.GetType())
	if !ok {
		return false
	}
	_, rType, ok := SplitResourceTypeName(revision.GetType())
	if !ok {
		return false
	}

	if !strings.HasSuffix(rType, "Revision") {
		return false
	}
	rType = strings.TrimSuffix(rType, "Revision")
	return pType == rType
}

// HasParent determines if the given resource has a parent resource or not
// based on the pattern(s) it defines having multiple resource ID segments
// or not. Incomplete or nil input returns false.
func HasParent(resource *annotations.ResourceDescriptor) bool {
	if resource == nil || len(resource.GetPattern()) == 0 {
		return false
	}

	for _, pattern := range resource.GetPattern() {
		// multiple ID variable segments indicates presence of parent
		if strings.Count(pattern, "{") > 1 {
			return true
		}
	}

	return false
}

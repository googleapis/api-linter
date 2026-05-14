// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"github.com/gertd/go-pluralize"
	"github.com/stoewer/go-strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var pluralizeClient = pluralize.NewClient()

// ToPlural converts a string to its plural form.
func ToPlural(s string) string {
	// Need to convert name to singular first to support none standard case such as persons, cactuses.
	// persons -> person -> people

	return pluralizeClient.Plural(pluralizeClient.Singular(s))
}

// ToSingular converts a string to its singular form.
func ToSingular(s string) string {
	return pluralizeClient.Singular(s)
}

// ResourceSingular returns the singular form of a resource name extracted from
// a batch method message name (e.g. "Books" from "BatchUpdateBooksRequest").
//
// It searches for a resource message with a google.api.resource annotation
// whose singular matches the expected singular of pluralName. This search
// covers messages in the same file as well as directly imported files, since
// the resource message is commonly defined in a separate file from the service
// and request/response messages.
//
// If a matching resource annotation is found, its singular is returned (in
// UpperCamelCase). Otherwise, it falls back to the go-pluralize library.
//
// This avoids incorrect singularization of words like "Metadata" (which
// go-pluralize converts to "Metadatum" using Latin grammar rules).
func ResourceSingular(pluralName string, m protoreflect.MessageDescriptor) string {
	if f := m.ParentFile(); f != nil {
		if s := findResourceSingularInFile(pluralName, f); s != "" {
			return s
		}
		// Also search directly imported files, since the resource message
		// is often defined in a separate file (e.g. impression_metadata.proto)
		// from the service file (impression_metadata_service.proto).
		imports := f.Imports()
		for i := 0; i < imports.Len(); i++ {
			if s := findResourceSingularInFile(pluralName, imports.Get(i).FileDescriptor); s != "" {
				return s
			}
		}
	}

	return pluralizeClient.Singular(pluralName)
}

// findResourceSingularInFile searches all messages in a file for a resource
// whose singular annotation matches the given pluralName.
func findResourceSingularInFile(pluralName string, f protoreflect.FileDescriptor) string {
	if f == nil {
		return ""
	}
	for i := 0; i < f.Messages().Len(); i++ {
		msg := f.Messages().Get(i)
		res := GetResource(msg)
		if res == nil {
			continue
		}
		s := GetResourceSingular(res)
		if s == "" {
			continue
		}
		// The singular from the annotation is typically lowerCamelCase
		// (e.g. "impressionMetadata"). Convert to UpperCamelCase for
		// comparison with the message name.
		upperSingular := strcase.UpperCamelCase(s)
		// Check if the pluralName matches: either the singular itself
		// (for uncountable nouns like "Metadata" where plural == singular)
		// or the go-pluralize plural of the singular.
		if pluralName == upperSingular || pluralName == pluralizeClient.Plural(upperSingular) {
			return upperSingular
		}
		// Also check if the message name matches the plural portion.
		if string(msg.Name()) == pluralName {
			return upperSingular
		}
	}
	return ""
}

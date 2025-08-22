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

package utils

import (
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// FindMessage looks for a message in a file and all imports within the
// same package.
func FindMessage(f protoreflect.FileDescriptor, name string) protoreflect.MessageDescriptor {
	// Default to using the current file's package.
	pkg := f.Package()

	// FileDescriptor.FindMessage requires fully-qualified message names;
	// attempt to infer that.
	if !strings.Contains(name, ".") {
		if pkg != "" {
			name = string(pkg) + "." + name
		}
	} else if !strings.HasPrefix(name, string(pkg)+".") {
		// If value is fully qualified, but from a different package,
		// accommodate that.
		pkg = protoreflect.FullName(name[:strings.LastIndex(name, ".")])
	}

	files := &protoregistry.Files{}
	RegisterFileRecursive(f, files)

	// Attempt to find the message in the file provided.
	if d, err := files.FindDescriptorByName(protoreflect.FullName(name)); err == nil {
		if m, ok := d.(protoreflect.MessageDescriptor); ok {
			// If the message's package is not what we expect, then it is
			// the wrong message.
			if m.ParentFile().Package() == pkg {
				return m
			}
		}
	}

	// Whelp, no luck. Too bad.
	return nil
}

// FindMethod searches a file and all imports within the same package, and
// returns the method with a given name, or nil if none is found.
func FindMethod(f protoreflect.FileDescriptor, name string) protoreflect.MethodDescriptor {
	for _, s := range getServices(f) {
		if m := s.Methods().ByName(protoreflect.Name(name)); m != nil {
			return m
		}
	}
	return nil
}

// getServices finds all services in a file and all imports within the
// same package.
func getServices(f protoreflect.FileDescriptor) []protoreflect.ServiceDescriptor {
	var answer []protoreflect.ServiceDescriptor
	for i := 0; i < f.Services().Len(); i++ {
		answer = append(answer, f.Services().Get(i))
	}
	for i := 0; i < f.Imports().Len(); i++ {
		dep := f.Imports().Get(i)
		if f.Package() == dep.Package() {
			answer = append(answer, getServices(dep.FileDescriptor)...)
		}
	}
	return answer
}

// GetAllDependencies returns all dependencies.
func GetAllDependencies(file protoreflect.FileDescriptor) map[string]protoreflect.FileDescriptor {
	answer := map[string]protoreflect.FileDescriptor{file.Path(): file}
	for i := 0; i < file.Imports().Len(); i++ {
		f := file.Imports().Get(i).FileDescriptor
		if _, found := answer[f.Path()]; !found {
			answer[f.Path()] = f
			for name, f2 := range GetAllDependencies(f) {
				answer[name] = f2
			}
		}
	}
	return answer
}

type fieldSorter []protoreflect.FieldDescriptor

// Len is part of sort.Interface.
func (f fieldSorter) Len() int {
	return len(f)
}

// Swap is part of sort.Interface.
func (f fieldSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Less is part of sort.Interface. Compare field number.
func (f fieldSorter) Less(i, j int) bool {
	return f[i].Number() < f[j].Number()
}

// GetRepeatedMessageFields returns all fields labeled `repeated` of type
// Message in the given message, sorted in field number order.
func GetRepeatedMessageFields(m protoreflect.MessageDescriptor) []protoreflect.FieldDescriptor {
	var fields fieldSorter

	// If an unresolable message is fed into this helper, return empty slice.
	if m == nil {
		return fields
	}

	for i := 0; i < m.Fields().Len(); i++ {
		f := m.Fields().Get(i)
		if f.IsList() && f.Kind() == protoreflect.MessageKind {
			fields = append(fields, f)
		}
	}

	sort.Sort(fields)

	return fields
}

// FindFieldDotNotation returns a field descriptor from a given message that
// corresponds to the dot separated path e.g. "book.name". If the path is
// unresolable the method returns nil. This is especially useful for resolving
// path variables in google.api.http and nested fields in
// google.api.method_signature annotations.
func FindFieldDotNotation(msg protoreflect.MessageDescriptor, ref string) protoreflect.FieldDescriptor {
	path := strings.Split(ref, ".")
	end := len(path) - 1
	for i, seg := range path {
		field := msg.Fields().ByName(protoreflect.Name(seg))
		if field == nil {
			return nil
		}

		if m := field.Message(); m != nil && i != end {
			msg = m
			continue
		}

		return field
	}

	return nil
}

// RegisterFileRecursive recursively registers a file and its dependencies
// into a protoregistry.Files.
func RegisterFileRecursive(f protoreflect.FileDescriptor, files *protoregistry.Files) {
	if _, err := files.FindFileByPath(f.Path()); err == nil {
		return // Already registered.
	}
	err := files.RegisterFile(f)
	if err != nil {
		// Ignore errors. This is a best-effort registration.
	}
	for i := 0; i < f.Imports().Len(); i++ {
		RegisterFileRecursive(f.Imports().Get(i).FileDescriptor, files)
	}
}

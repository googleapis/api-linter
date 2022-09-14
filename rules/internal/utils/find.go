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

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FindMessage looks for a message in a file and all imports within the
// same package.
func FindMessage(f *desc.FileDescriptor, name string) *desc.MessageDescriptor {
	// FileDescriptor.FindMessage requires fully-qualified message names;
	// attempt to infer that.
	if !strings.Contains(name, ".") && f.GetPackage() != "" {
		name = f.GetPackage() + "." + name
	}

	// Attempt to find the message in the file provided.
	if m := f.FindMessage(name); m != nil {
		return m
	}

	// Attempt to find the message in any dependency files if they are in the
	// same package.
	for _, dep := range f.GetDependencies() {
		if f.GetPackage() == dep.GetPackage() {
			if m := FindMessage(dep, name); m != nil {
				return m
			}
		}
	}

	// Whelp, no luck. Too bad.
	return nil
}

// FindMethod searches a file and all imports within the same package, and
// returns the method with a given name, or nil if none is found.
func FindMethod(f *desc.FileDescriptor, name string) *desc.MethodDescriptor {
	for _, s := range getServices(f) {
		for _, m := range s.GetMethods() {
			if m.GetName() == name {
				return m
			}
		}
	}
	return nil
}

// getServices finds all services in a file and all imports within the
// same package.
func getServices(f *desc.FileDescriptor) []*desc.ServiceDescriptor {
	answer := f.GetServices()
	for _, dep := range f.GetDependencies() {
		if f.GetPackage() == dep.GetPackage() {
			answer = append(answer, getServices(dep)...)
		}
	}
	return answer
}

// GetAllDependencies returns all dependencies.
func GetAllDependencies(file *desc.FileDescriptor) map[string]*desc.FileDescriptor {
	answer := map[string]*desc.FileDescriptor{file.GetName(): file}
	for _, f := range file.GetDependencies() {
		if _, found := answer[f.GetName()]; !found {
			answer[f.GetName()] = f
			for name, f2 := range GetAllDependencies(f) {
				answer[name] = f2
			}
		}
	}
	return answer
}

type fieldSorter []*desc.FieldDescriptor

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
	return f[i].GetNumber() < f[j].GetNumber()
}

// GetRepeatedMessageFields returns all fields labeled `repeated` of type
// Message in the given message, sorted in field number order.
func GetRepeatedMessageFields(m *desc.MessageDescriptor) []*desc.FieldDescriptor {
	var fields fieldSorter

	// If an unresolable message is fed into this helper, return empty slice.
	if m == nil {
		return fields
	}

	for _, f := range m.GetFields() {
		if f.IsRepeated() && f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			fields = append(fields, f)
		}
	}

	sort.Sort(fields)

	return fields
}

// FindFieldDotNotation attempts to find the field within the given message
// identified by the dot-notation path e.g. book.name. This is especially useful
// for resolving path variables in google.api.http and nested
// fields in google.api.method_signature annotations.
func FindFieldDotNotation(msg *desc.MessageDescriptor, f string) *desc.FieldDescriptor {
	path := strings.Split(f, ".")
	for _, seg := range path {
		field := msg.FindFieldByName(seg)
		if field == nil {
			return nil
		}

		if m := field.GetMessageType(); m != nil {
			msg = m
			continue
		}

		return field
	}

	return nil
}

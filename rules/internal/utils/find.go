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
	"strings"

	"github.com/jhump/protoreflect/desc"
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

// Copyright 2019 Google LLC
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

package lint

import (
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// SyntaxLocation returns the location of the syntax definition in a file
// descriptor.
//
// If the location can not be found (for example, because there is no syntax
// statement), it returns beginning of file.
func SyntaxLocation(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return sourceInfoRegistry.sourceInfo(f).findLocation([]int32{12}) // syntax == 12
}

// PackageLocation returns the location of the package definition in a file
// descriptor.
//
// If the location can not be found (for example, because there is no package
// statement), it returns beginning of file.
func PackageLocation(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return sourceInfoRegistry.sourceInfo(f).findLocation([]int32{2}) // package == 2
}

// DescriptorNameLocation returns the precise location for a descriptor's name.
func DescriptorNameLocation(d desc.Descriptor) *dpb.SourceCodeInfo_Location {
	// All descriptors seem to have `string name = 1`, so this conveniently works.
	path := append(d.GetSourceInfo().Path, 1)
	return sourceInfoRegistry.sourceInfo(d.GetFile()).findLocation(path)
}

type sourceInfo map[string]*dpb.SourceCodeInfo_Location

// findLocation returns the best location it can find for the given path.
//
// If the requested path can not be found, it climbs the ancestry tree until
// it finds one, and ultimately returns a location corresponding to the
// beginning of the file if it can not find anything.
func (si sourceInfo) findLocation(path []int32) *dpb.SourceCodeInfo_Location {
	// Base case: If we have no path, return nil.
	if len(path) == 0 {
		return nil
	}

	// If the path exists in the source info registry, return that object.
	if loc, ok := si[strPath(path)]; ok {
		return loc
	}

	// We could not find the path; return nil.
	return nil
}

// The source map registry is a singleton that computes a source map for
// any file descriptor that it is given, but then caches it to avoid computing
// the source map for the same file descriptors over and over.
type sourceInfoRegistryType map[*desc.FileDescriptor]sourceInfo

// Each location has a path defined as an []int32, but we can not
// use slices as keys, so compile them into a string.
func strPath(segments []int32) (p string) {
	for i, segment := range segments {
		if i > 0 {
			p += ","
		}
		p += string(segment)
	}
	return
}

// sourceInfo compiles the source info object for a given file descriptor.
// It also caches this into a registry, so subsequent calls using the same
// descriptor will return the same object.
func (sir sourceInfoRegistryType) sourceInfo(fd *desc.FileDescriptor) sourceInfo {
	answer, ok := sir[fd]
	if !ok {
		answer = sourceInfo{}

		// This file descriptor does not yet have a source info map.
		// Compile one.
		for _, loc := range fd.AsFileDescriptorProto().GetSourceCodeInfo().GetLocation() {
			answer[strPath(loc.Path)] = loc
		}

		// Now that we calculated all of this, cache it on the registry so it
		// does not need to be calculated again.
		sir[fd] = answer
	}
	return answer
}

var sourceInfoRegistry = sourceInfoRegistryType{}

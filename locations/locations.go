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

// Package locations provides functions to get the location of a particular part
// of a descriptor, allowing Problems to be attached to just a descriptor's
// name, type, etc.. This allows for better auto-replacement functionality in
// code review tools.
//
// All functions in this package accept a descriptor and return a
// protobuf SourceCodeInfo_Location object, which can be passed directly
// to the Location property on Problem.
package locations

import (
	"github.com/jhump/protoreflect/desc"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// pathLocation returns the precise location for a given descriptor and path.
// It combines the path of the descriptor itself with any path provided appended.
func pathLocation(d desc.Descriptor, path ...int) *dpb.SourceCodeInfo_Location {
	fullPath := d.GetSourceInfo().GetPath()
	for _, i := range path {
		fullPath = append(fullPath, int32(i))
	}
	return sourceInfoRegistry.sourceInfo(d.GetFile()).findLocation(fullPath)
}

type sourceInfo map[string]*dpb.SourceCodeInfo_Location

// findLocation returns the Location for a given path.
func (si sourceInfo) findLocation(path []int32) *dpb.SourceCodeInfo_Location {
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

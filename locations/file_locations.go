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

package locations

import (
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// FileSyntax returns the location of the syntax definition in a file descriptor.
//
// If the location can not be found (for example, because there is no syntax
// statement), it returns nil.
func FileSyntax(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 12) // FileDescriptor.syntax == 12
}

// FilePackage returns the location of the package definition in a file descriptor.
//
// If the location can not be found (for example, because there is no package
// statement), it returns nil.
func FilePackage(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 2) // FileDescriptor.package == 2
}

// FileCsharpNamespace returns the location of the csharp_namespace file option
// in a file descriptor.
//
// If the location can not be found (for example, because there is no
// csharp_namespace option), it returns nil.
func FileCsharpNamespace(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, 37) // 8 == options, 37 == csharp_namespace
}

// FileJavaPackage returns the location of the java_package file option
// in a file descriptor.
//
// If the location can not be found (for example, because there is no
// java_package option), it returns nil.
func FileJavaPackage(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, 1) // 8 == options, 1 == java_package
}

// FilePhpNamespace returns the location of the php_namespace file option
// in a file descriptor.
//
// If the location can not be found (for example, because there is no
// php_namespace option), it returns nil.
func FilePhpNamespace(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, 41) // 8 == options, 41 == php_namespace
}

// FileRubyPackage returns the location of the ruby_package file option
// in a file descriptor.
//
// If the location can not be found (for example, because there is no
// ruby_package option), it returns nil.
func FileRubyPackage(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, 45) // 8 == options, 45 == ruby_package
}

// FileResourceDefinition returns the precise location of the `google.api.resource_definition`
// annotation.
func FileResourceDefinition(f *desc.FileDescriptor, index int) *dpb.SourceCodeInfo_Location {
	// 8 == options
	return pathLocation(f, 8, int(apb.E_ResourceDefinition.TypeDescriptor().Number()), index)
}

// FileImport returns the location of the import on the given `index`, or `nil`
// if no import with such `index` is found.
func FileImport(f *desc.FileDescriptor, index int) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 3, index) // 3 == dependency
}

// FileCCEnableArenas returns the location of the `cc_enable_arenas` option.
func FileCCEnableArenas(f *desc.FileDescriptor) *dpb.SourceCodeInfo_Location {
	return pathLocation(f, 8, 31) // 8 == (file) options, 31 == cc_enable_arenas
}

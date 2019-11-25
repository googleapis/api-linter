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

// Descriptor defines a descriptor in a Protobuf file.
type Descriptor interface {
	SourceInfo() SourceInfo
}

// SourceInfo defines source information about a descriptor.
type SourceInfo interface {
	File() FileInfo
	LeadingComments() string
}

// FileInfo defines source information about a file.
type FileInfo interface {
	Path() string
	Comments() string
}

// Rule defines a lint rule that checks Google Protobuf APIs.
//
// Anything that satisfies this interface can be used as a rule,
// but most rule authors will want to use the implementations provided.
type Rule interface {
	// Name returns the name of the rule.
	Name() RuleName

	// Lint accepts a Descriptor and lints it,
	// returning a slice of Problem objects it finds.
	Lint(Descriptor) []Problem
}

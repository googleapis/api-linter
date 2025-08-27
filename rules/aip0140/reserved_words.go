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

package aip0140

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/locations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var reservedWords = &lint.FieldRule{
	Name: lint.NewRuleName(140, "reserved-words"),
	LintField: func(f protoreflect.FieldDescriptor) []lint.Problem {
		if name := f.Name(); reservedWordsSet.Contains(string(name)) {
			return []lint.Problem{{
				Message:    fmt.Sprintf("%q is a reserved word in a common language and should not be used.", name),
				Descriptor: f,
				Location:   locations.DescriptorName(f),
			}}
		}
		return nil
	},
}

var reservedWordsSet = stringset.New(
	"abstract",     // Java, JavaScript
	"and",          // Python
	"arguments",    // JavaScript
	"as",           // Python
	"assert",       // Java, Python
	"async",        // Python
	"await",        // JavaScript, Python
	"boolean",      // Java, JavaScript
	"break",        // Java, JavaScript, Python
	"byte",         // Java, JavaScript
	"case",         // Java, JavaScript
	"catch",        // Java, JavaScript
	"char",         // Java, JavaScript
	"class",        // Java, JavaScript, Python
	"const",        // Java, JavaScript
	"continue",     // Java, JavaScript, Python
	"debugger",     // JavaScript
	"def",          // Python
	"default",      // Java, JavaScript
	"del",          // Python
	"delete",       // JavaScript
	"do",           // Java, JavaScript
	"double",       // Java, JavaScript
	"elif",         // Python
	"else",         // Java, JavaScript, Python
	"enum",         // Java, JavaScript
	"eval",         // JavaScript
	"except",       // Python
	"export",       // JavaScript
	"extends",      // Java, JavaScript
	"false",        // JavaScript
	"final",        // Java, JavaScript
	"finally",      // Java, JavaScript, Python
	"float",        // Java, JavaScript
	"for",          // Java, JavaScript, Python
	"from",         // Python
	"function",     // JavaScript
	"global",       // Python
	"goto",         // Java, JavaScript
	"if",           // Java, JavaScript, Python
	"implements",   // Java, JavaScript
	"import",       // Java, JavaScript, Python
	"in",           // JavaScript, JavaScript, Python
	"instanceof",   // Java, JavaScript
	"int",          // Java, JavaScript
	"interface",    // Java, JavaScript
	"is",           // Python
	"lambda",       // Python
	"let",          // JavaScript
	"long",         // Java, JavaScript
	"native",       // Java, JavaScript
	"new",          // Java, JavaScript
	"nonlocal",     // Python
	"not",          // Python
	"null",         // JavaScript
	"or",           // Python
	"package",      // Java, JavaScript
	"pass",         // Python
	"private",      // Java, JavaScript
	"protected",    // Java, JavaScript
	"public",       // Java, JavaScript
	"raise",        // Python
	"return",       // Java, JavaScript, Python
	"short",        // Java, JavaScript
	"static",       // Java, JavaScript
	"strictfp",     // Java
	"super",        // Java, JavaScript
	"switch",       // Java, JavaScript
	"synchronized", // Java, JavaScript
	"this",         // Java, JavaScript
	"throw",        // Java, JavaScript
	"throws",       // Java, JavaScript
	"transient",    // Java, JavaScript
	"true",         // JavaScript
	"try",          // Java, JavaScript, Python
	"typeof",       // JavaScript
	"var",          // JavaScript
	"void",         // Java, JavaScript
	"volatile",     // Java, JavaScript
	"while",        // Java, JavaScript, Python
	"with",         // JavaScript, Python
	"yield",        // JavaScript, Python
)

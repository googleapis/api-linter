// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0215

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	"github.com/googleapis/api-linter/v2/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var foreignTypeReference = &lint.FieldRule{
	Name: lint.NewRuleName(215, "foreign-type-reference"),
	OnlyIf: func(fd protoreflect.FieldDescriptor) bool {
		return fd.Kind() == protoreflect.MessageKind
	},
	LintField: func(fd protoreflect.FieldDescriptor) []lint.Problem {
		curPkg := getNormalizedPackage(fd)
		if curPkg == "" {
			return nil // Empty or unavailable package.
		}
		msg := fd.Message()
		if msg == nil {
			return nil // Couldn't resolve type.
		}
		msgPkg := getNormalizedPackage(msg)
		if msgPkg == "" {
			return nil // Empty or unavailable package.
		}

		if utils.IsCommonProto(msg.ParentFile()) {
			return nil // reference to a well known proto package.
		}

		if strings.HasSuffix(msgPkg, ".type") {
			return nil // AIP-213 component type.
		}

		if curPkg != msgPkg {
			return []lint.Problem{{
				Message:    fmt.Sprintf("foreign type referenced, current field in %q message in %q", curPkg, msgPkg),
				Descriptor: fd,
			}}
		}

		return nil
	},
}

// Regexp to capture everything up to a versioned segment.
var versionedPrefix = regexp.MustCompile(`^.*\.v[\d]+(p[\d]+)?(alpha|beta|eap|test)?[\d]*`)

// getNormalizedPackage returns a normalized package path.
// If package cannot be resolved it returns the empty string.
// If the package path has a "versioned" segment, the path is truncated to that segment.
func getNormalizedPackage(d protoreflect.Descriptor) string {
	f := d.ParentFile()
	if f == nil {
		return ""
	}
	pkg := string(f.Package())
	if normPkg := versionedPrefix.FindString(pkg); normPkg != "" {
		pkg = normPkg
	}
	return pkg
}

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

package aip0191

import (
	"github.com/commure/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var fileLayout = &lint.FileRule{
	Name: lint.NewRuleName(191, "file-layout"),
	LintFile: func(f *desc.FileDescriptor) (problems []lint.Problem) {
		// Verify that services precede messages.
		if len(f.GetMessageTypes()) > 0 {
			firstMessage := f.GetMessageTypes()[0]
			for _, service := range f.GetServices() {
				if isAfter(firstMessage, service) {
					problems = append(problems, lint.Problem{
						Message:    "Services should precede all messages.",
						Descriptor: service,
					})
				}
			}
		}

		// Verify that messages precede top-level enums.
		if len(f.GetEnumTypes()) > 0 {
			firstEnum := f.GetEnumTypes()[0]
			for _, message := range f.GetMessageTypes() {
				if isBefore(message, firstEnum) {
					problems = append(problems, lint.Problem{
						Message:    "Messages should precede all top-level enums.",
						Descriptor: firstEnum,
					})
					break // Sending this over and over would be obnoxious.
				}
			}
		}

		return
	},
}

// isBefore returns true if `d` is known to precede `anchor` in the file.
//
// NOTE: A false value here may indicate that there is no source info at all;
//       use `isAfter` if the goal is to know that `d` comes after `anchor`.
func isBefore(anchor desc.Descriptor, d desc.Descriptor) bool {
	return d.GetSourceInfo().GetSpan()[0] < anchor.GetSourceInfo().GetSpan()[0]
}

// isBefore returns true if `d` is known to follow `anchor` in the file.
//
// NOTE: A false value here may indicate that there is no source info at all;
//       use `isBefore` if the goal is to know that `d` comes before `anchor`.
func isAfter(anchor desc.Descriptor, d desc.Descriptor) bool {
	return d.GetSourceInfo().GetSpan()[0] > anchor.GetSourceInfo().GetSpan()[0]
}

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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrPathNotFound is the returned error when a path is not found.
	ErrPathNotFound = errors.New("source: path not found")
	// ErrSourceInfoNotAvailable is the returned error when creating a source
	// but the source information is not available.
	ErrSourceInfoNotAvailable = errors.New("source: source information is not available")
)

// Comments describes a collection of comments associate with an element,
// which contains leading, trailing, and leading-detached comments, in a
// source code file.
type Comments struct {
	LeadingComments         string
	TrailingComments        string
	LeadingDetachedComments []string
}

const sep = ","

// locPath represents a path in the SourceCodeInfo_Location,
// and this serves as a map key.
// It's a string representation of a slice because slices
// cannot be map keys.
// Representation: integers separated by commas. No spaces.
// Example: [4, 3, 2, 7] --> "4,3,2,7"
// See descriptor.proto for more explanation of semantics.
type locPath string

// newLocPath return a locPath from a list of index.
func newLocPath(p ...int) locPath {
	a := []string{}
	for _, i := range p {
		a = append(a, strconv.Itoa(i))
	}
	return locPath(strings.Join(a, sep))
}

// buildLocPathMap creates a map of locPath to *descriptorpb.SourceCodeInfo_Location
// from *descriptorpb.SourceCodeInfo.
func buildLocPathMap(sci protoreflect.SourceLocations) map[locPath]*protoreflect.SourceLocation {
	m := make(map[locPath]*protoreflect.SourceLocation)
	if sci == nil {
		return m
	}

	for i := 0; i < sci.Len(); i++ {
		loc := sci.Get(i)
		path := make([]int, len(loc.Path))
		for k := range path {
			path[k] = int(loc.Path[k])
		}
		m[newLocPath(path...)] = &loc
	}
	return m
}

// DescriptorSource represents a map of locPath to *descriptorpb.SourceCodeInfo_Location.
type DescriptorSource struct {
	m map[locPath]*protoreflect.SourceLocation
}

// newDescriptorSource creates a new DescriptorSource from a FileDescriptor.
// If source code information is not available, returns (nil, ErrSourceInfoNotAvailable).
func newDescriptorSource(f protoreflect.FileDescriptor) (DescriptorSource, error) {
	if f.SourceLocations() == nil {
		return DescriptorSource{}, ErrSourceInfoNotAvailable
	}
	return DescriptorSource{m: buildLocPathMap(f.SourceLocations())}, nil
}

// findLocationByPath returns a `Location` if found in the map,
// and (nil, ErrPathNotFound) if not found.
func (s DescriptorSource) findLocationByPath(path []int) (Location, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Location{}, ErrPathNotFound
	}
	return extractLocation(l), nil
}

// findCommentsByPath returns a `Comments` for the path. If not found, returns
// (nil, ErrCommentsNotFound).
func (s DescriptorSource) findCommentsByPath(path []int) (Comments, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Comments{}, ErrPathNotFound
	}
	return Comments{
		LeadingComments:         l.LeadingComments,
		TrailingComments:        l.TrailingComments,
		LeadingDetachedComments: l.LeadingDetachedComments,
	}, nil
}

func extractLocation(loc *protoreflect.SourceLocation) Location {
	return Location{
		Start: Position{
			Line:   loc.StartLine + 1,
			Column: loc.StartColumn + 1,
		},
		End: Position{
			Line:   loc.EndLine + 1,
			Column: loc.EndColumn + 1,
		},
	}
}

// SyntaxLocation returns the location of the syntax definition.
func (s DescriptorSource) SyntaxLocation() (Location, error) {
	return s.findLocationByPath([]int{syntaxTag})
}

// SyntaxComments returns the comments of the syntax definition.
func (s DescriptorSource) SyntaxComments() (Comments, error) {
	return s.findCommentsByPath([]int{syntaxTag})
}

// PackageLocation returns the location of the package definition.
func (s DescriptorSource) PackageLocation() (Location, error) {
	return s.findLocationByPath([]int{packageTag})
}

// PackageComments returns the comments of the package definition.
func (s DescriptorSource) PackageComments() (Comments, error) {
	return s.findCommentsByPath([]int{packageTag})
}

// DescriptorLocation returns a `Location` for the given descriptor.
// If not found, returns (nil, ErrPathNotFound).
func (s DescriptorSource) DescriptorLocation(d protoreflect.Descriptor) (Location, error) {
	return s.findLocationByPath(getPath(d))
}

// DescriptorLocationOrFileStart returns a `Location` for the given descriptor. If there was an
// error finding the location, it returns the start location of the file instead (that is,
// line 1 column 1 as the Start and End Positions).
func (s DescriptorSource) DescriptorLocationOrFileStart(d protoreflect.Descriptor) Location {
	loc, err := s.DescriptorLocation(d)
	if err != nil {
		return Location{
			Start: Position{Line: 1, Column: 1},
			End:   Position{Line: 1, Column: 1},
		}
	}
	return loc
}

func getPath(d protoreflect.Descriptor) []int {
	path := []int{}
	for p := d; p != nil && !isFileDescriptor(p); p = p.Parent() {
		path = append(path, p.Index(), getDescriptorTag(p))
	}
	reverseInts(path)
	return path
}

const syntaxTag = 12
const packageTag = 2

var enumTagInFile = 5
var enumTagInMessage = 4
var enumValueTag = 2
var fieldTag = 2
var extensionTagInFile = 7
var extensionTagInMessage = 6
var messageTagInFile = 4
var nestedMessageTag = 3
var oneofTag = 8
var serviceTag = 6
var methodTag = 2

func getDescriptorTag(d protoreflect.Descriptor) int {
	switch d.(type) {
	case protoreflect.EnumDescriptor:
		if isTopLevelDescriptor(d) {
			return enumTagInFile
		}
		return enumTagInMessage
	case protoreflect.EnumValueDescriptor:
		return enumValueTag
	case protoreflect.FieldDescriptor:
		if isFieldExtension(d) {
			if isTopLevelDescriptor(d) {
				return extensionTagInFile
			}
			return extensionTagInMessage
		}
		return fieldTag
	case protoreflect.MessageDescriptor:
		if isTopLevelDescriptor(d) {
			return messageTagInFile
		}
		return nestedMessageTag
	case protoreflect.MethodDescriptor:
		return methodTag
	case protoreflect.OneofDescriptor:
		return oneofTag
	case protoreflect.ServiceDescriptor:
		return serviceTag
	default:
		return 0
	}
}

func isFieldExtension(d protoreflect.Descriptor) bool {
	f, ok := d.(protoreflect.FieldDescriptor)
	return ok && f.IsExtension()
}

func isFileDescriptor(d protoreflect.Descriptor) bool {
	_, ok := d.(protoreflect.FileDescriptor)
	return ok
}

func isTopLevelDescriptor(d protoreflect.Descriptor) bool {
	p := d.Parent()
	_, ok := p.(protoreflect.FileDescriptor)
	return ok
}

// DescriptorComments returns a `Comments` for the given descriptor.
// If not found, returns (nil, ErrCommentsNotFound).
func (s DescriptorSource) DescriptorComments(d protoreflect.Descriptor) (Comments, error) {
	return s.findCommentsByPath(getPath(d))
}

func reverseInts(a []int) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

// isRuleDisabled check if a rule is disabled for a descriptor
// in the comments.
func (s DescriptorSource) isRuleDisabled(name RuleName, d protoreflect.Descriptor) bool {
	commentsToCheck := s.fileComments().LeadingDetachedComments

	for ; d != nil; d = d.Parent() {
		comments, err := s.DescriptorComments(d)

		if err != nil {
			continue
		}

		commentsToCheck = append(commentsToCheck, comments.LeadingComments, comments.TrailingComments)
	}

	return stringsContains(commentsToCheck, ruleDisablingComment(name))
}

// isRuleDisabledInFile checks the proto file comments only to see if a rule named name is disabled.
func (s DescriptorSource) isRuleDisabledInFile(name RuleName) bool {
	return s.isRuleDisabled(name, nil)
}

func stringsContains(comments []string, s string) bool {
	for _, c := range comments {
		if strings.Contains(c, s) {
			return true
		}
	}
	return false
}

func ruleDisablingComment(name RuleName) string {
	return fmt.Sprintf("(-- api-linter: %s=disabled --)", name)
}

func (s DescriptorSource) fileComments() Comments {
	comments, _ := s.SyntaxComments()
	return comments
}

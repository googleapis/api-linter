package lint

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

var (
	// ErrLocationNotFound is the returned error when a location is not found.
	ErrLocationNotFound = errors.New("source: location not found")
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
// See descriptor.proto for more explantion of semantics.
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
func buildLocPathMap(sci *descriptorpb.SourceCodeInfo) map[locPath]*descriptorpb.SourceCodeInfo_Location {
	m := make(map[locPath]*descriptorpb.SourceCodeInfo_Location)
	if sci == nil {
		return m
	}

	for _, loc := range sci.GetLocation() {
		var path []int
		for _, v := range loc.GetPath() {
			path = append(path, int(v))
		}
		m[newLocPath(path...)] = loc
	}
	return m
}

// DescriptorSource represents a map of locPath to *descriptorpb.SourceCodeInfo_Location.
type DescriptorSource struct {
	m map[locPath]*descriptorpb.SourceCodeInfo_Location
}

// NewDescriptorSource creates a new DescriptorSource from a FileDescriptorProto.
// If source code information is not available, returns (nil, ErrSourceInfoNotAvailable).
func NewDescriptorSource(f *descriptorpb.FileDescriptorProto) (DescriptorSource, error) {
	if f.GetSourceCodeInfo() == nil {
		return DescriptorSource{}, ErrSourceInfoNotAvailable
	}
	return DescriptorSource{m: buildLocPathMap(f.GetSourceCodeInfo())}, nil
}

// findLocationByPath returns a `Location` if found in the map,
// and (nil, ErrLocationNotFound) if not found.
func (s DescriptorSource) findLocationByPath(path []int) (Location, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Location{}, ErrLocationNotFound
	}
	return newLocationFromSpan(l.GetSpan()), nil
}

// findCommentsByPath returns a `Comments` for the path. If not found, returns
// (nil, ErrCommentsNotFound).
func (s DescriptorSource) findCommentsByPath(path []int) (Comments, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Comments{}, ErrLocationNotFound
	}
	return Comments{
		LeadingComments:         l.GetLeadingComments(),
		TrailingComments:        l.GetTrailingComments(),
		LeadingDetachedComments: l.GetLeadingDetachedComments(),
	}, nil
}

func newLocationFromSpan(span []int32) Location {
	if len(span) == 4 {
		return Location{
			Start: Position{
				Line:   int(span[0]),
				Column: int(span[1]),
			},
			End: Position{
				Line:   int(span[2]),
				Column: int(span[3]),
			},
		}
	}

	if len(span) == 3 {
		return Location{
			Start: Position{
				Line:   int(span[0]),
				Column: int(span[1]),
			},
			End: Position{
				Line:   int(span[0]),
				Column: int(span[2]),
			},
		}
	}

	return Location{}
}

// SyntaxLocation returns the location of the syntax definition.
func (s DescriptorSource) SyntaxLocation() (Location, error) {
	return s.findLocationByPath([]int{syntaxTag})
}

// SyntaxComments returns the comments of the syntax definition.
func (s DescriptorSource) SyntaxComments() (Comments, error) {
	return s.findCommentsByPath([]int{syntaxTag})
}

// DescriptorLocation returns a `Location` for the given descriptor.
// If not found, returns (nil, ErrLocationNotFound).
func (s DescriptorSource) DescriptorLocation(d protoreflect.Descriptor) (Location, error) {
	return s.findLocationByPath(getPath(d))
}

func getPath(d protoreflect.Descriptor) []int {
	path := []int{}
	for p := d; !isFileDescriptor(p); p, _ = p.Parent() {
		path = append(path, p.Index(), getDescriptorTag(p))
	}
	reverseInts(path)
	return path
}

const syntaxTag = 12

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
	return ok && f.ExtendedType() != nil
}

func isFileDescriptor(d protoreflect.Descriptor) bool {
	_, ok := d.(protoreflect.FileDescriptor)
	return ok
}

func isTopLevelDescriptor(d protoreflect.Descriptor) bool {
	p, _ := d.Parent()
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

// IsRuleDisabled check if a rule is disabled for a descriptor
// in the comments.
func (s DescriptorSource) IsRuleDisabled(id RuleID, d protoreflect.Descriptor) bool {
	comments, err := s.DescriptorComments(d)
	if err != nil {
		return true
	}

	commentsToCheck := []string{
		comments.LeadingComments,
		comments.TrailingComments,
	}
	commentsToCheck = append(commentsToCheck, s.fileComments().LeadingDetachedComments...)

	return stringsContains(commentsToCheck, ruleDisablingComment(id))
}

func stringsContains(comments []string, s string) bool {
	for _, c := range comments {
		if strings.Contains(c, s) {
			return true
		}
	}
	return false
}

func ruleDisablingComment(id RuleID) string {
	name := id.Set + "." + id.Name
	if id.Set == "" || id.Set == "core" {
		name = id.Name
	}
	return "(-- api-linter: " + name + "=disabled --)"
}

func (s DescriptorSource) fileComments() Comments {
	comments, _ := s.SyntaxComments()
	return comments
}

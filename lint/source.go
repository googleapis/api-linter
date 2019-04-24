package lint

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/lint/protowalk"
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

// newDescriptorSource creates a new DescriptorSource from a FileDescriptorProto.
// If source code information is not available, returns (nil, ErrSourceInfoNotAvailable).
func newDescriptorSource(f *descriptorpb.FileDescriptorProto) (DescriptorSource, error) {
	if f.GetSourceCodeInfo() == nil {
		return DescriptorSource{}, ErrSourceInfoNotAvailable
	}
	return DescriptorSource{m: buildLocPathMap(f.GetSourceCodeInfo())}, nil
}

// findLocationByPath returns a `Location` if found in the map,
// and (nil, ErrPathNotFound) if not found.
func (s DescriptorSource) findLocationByPath(path []int) (*Location, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return nil, ErrPathNotFound
	}
	return newLocationFromSpan(l.GetSpan())
}

// findCommentsByPath returns a `Comments` for the path. If not found, returns
// (nil, ErrCommentsNotFound).
func (s DescriptorSource) findCommentsByPath(path []int) (Comments, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Comments{}, ErrPathNotFound
	}
	return Comments{
		LeadingComments:         l.GetLeadingComments(),
		TrailingComments:        l.GetTrailingComments(),
		LeadingDetachedComments: l.GetLeadingDetachedComments(),
	}, nil
}

func newLocationFromSpan(span []int32) (*Location, error) {
	if len(span) == 4 {
		start := NewPosition(int(span[0]), int(span[1]))
		end := NewPosition(int(span[2]), int(span[3]))
		return NewLocation(start, end), nil
	}

	if len(span) == 3 {
		start := NewPosition(int(span[0]), int(span[1]))
		end := NewPosition(int(span[0]), int(span[2]))
		return NewLocation(start, end), nil
	}

	return nil, fmt.Errorf("source: %v is not a valid span to create a Location", span)
}

// SyntaxLocation returns the location of the syntax definition.
func (s DescriptorSource) SyntaxLocation() (*Location, error) {
	return s.findLocationByPath([]int{syntaxTag})
}

// SyntaxComments returns the comments of the syntax definition.
func (s DescriptorSource) SyntaxComments() (Comments, error) {
	return s.findCommentsByPath([]int{syntaxTag})
}

// DescriptorLocation returns a `Location` for the given descriptor.
// If not found, returns (nil, ErrPathNotFound).
func (s DescriptorSource) DescriptorLocation(d protoreflect.Descriptor) (*Location, error) {
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

var ruleDisabledPattern = fmt.Sprintf("\\(-- api-linter: (%s)=disabled --\\)", ruleNamePattern)
var ruleDisabledRegexp = regexp.MustCompile(ruleDisabledPattern)

func findDisabledRules(s string) []string {
	match := ruleDisabledRegexp.FindAllStringSubmatch(s, -1)
	results := make([]string, len(match))
	for i, m := range match {
		results[i] = m[1]
	}
	return results
}

type disabledRuleFinder struct {
	ruleLocs map[string][]*Location
	source   DescriptorSource
}

func (finder *disabledRuleFinder) isRuleDisabledAtLocation(name RuleName, loc *Location) bool {
	for _, l := range finder.ruleLocs[string(name)] {
		if l.contains(loc) {
			return true
		}
	}
	return false
}

func (finder *disabledRuleFinder) isRuleDisabledAtDescriptor(name RuleName, d protoreflect.Descriptor) bool {
	loc, err := finder.source.DescriptorLocation(d)
	if err != nil {
		return false
	}
	return finder.isRuleDisabledAtLocation(name, loc)
}

const uintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

const (
	maxInt = 1<<(uintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
)

func (finder *disabledRuleFinder) Consume(d protoreflect.Descriptor) error {
	comments, _ := finder.source.DescriptorComments(d)
	location, _ := finder.source.DescriptorLocation(d)
	s := comments.LeadingComments + "\n" + comments.TrailingComments
	if _, ok := d.(protoreflect.FileDescriptor); ok {
		comments, _ = finder.source.SyntaxComments()
		location = NewLocation(NewPosition(0, 0), NewPosition(maxInt, maxInt))
		s = strings.Join(comments.LeadingDetachedComments, "")
	}

	rules := findDisabledRules(s)
	for _, rl := range rules {
		finder.ruleLocs[rl] = append(finder.ruleLocs[rl], location)
	}
	return nil
}

func newDisabledRuleFinder(f protoreflect.FileDescriptor, s DescriptorSource) *disabledRuleFinder {
	finder := &disabledRuleFinder{
		ruleLocs: make(map[string][]*Location),
		source:   s,
	}
	protowalk.Walk(f, finder)
	return finder
}

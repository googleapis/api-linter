package lint

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang/protobuf/reflect/protoreflect"
	dpb "github.com/golang/protobuf/types/descriptor"
)

var (
	// ErrLocationNotFound is the returned error when a location is not found.
	ErrLocationNotFound = errors.New("source: location not found")
	// ErrCommentsNotFound is the returned error when a comments is not found.
	ErrCommentsNotFound = errors.New("source: comments not found")
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

// buildLocPathMap creates a map of locPath to *dpb.SourceCodeInfo_Location
// from *dpb.SourceCodeInfo.
func buildLocPathMap(sci *dpb.SourceCodeInfo) map[locPath]*dpb.SourceCodeInfo_Location {
	m := make(map[locPath]*dpb.SourceCodeInfo_Location)
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

// DescriptorSource represents a map of locPath to *dpb.SourceCodeInfo_Location.
type DescriptorSource struct {
	m map[locPath]*dpb.SourceCodeInfo_Location
}

// NewDescriptorSource creates a new DescriptorSource from a FileDescriptorProto.
// If source code information is not available, returns (nil, ErrSourceInfoNotAvailable).
func NewDescriptorSource(f *dpb.FileDescriptorProto) (*DescriptorSource, error) {
	if f.GetSourceCodeInfo() == nil {
		return nil, ErrSourceInfoNotAvailable
	}
	return &DescriptorSource{m: buildLocPathMap(f.GetSourceCodeInfo())}, nil
}

// FindLocationByPath returns a `Location` if found in the map,
// and (nil, ErrLocationNotFound) if not found.
func (s DescriptorSource) FindLocationByPath(path []int) (Location, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Location{}, ErrLocationNotFound
	}
	return newLocationFromSpan(l.GetSpan()), nil
}

// FindCommentsByPath returns a `Comments` for the path. If not found, returns
// (nil, ErrCommentsNotFound).
func (s DescriptorSource) FindCommentsByPath(path []int) (Comments, error) {
	l := s.m[newLocPath(path...)]
	if l == nil {
		return Comments{}, ErrCommentsNotFound
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

// FindLocationByDescriptor returns a `Location` for the given descriptor.
// If not found, returns (nil, ErrLocationNotFound).
func (s DescriptorSource) FindLocationByDescriptor(d protoreflect.Descriptor) (Location, error) {
	return s.FindLocationByPath(getPath(d))
}

func getPath(d protoreflect.Descriptor) []int {
	if d.Index() == 0 {
		return []int{}
	}
	path := []int{d.Index()}
	p, found := d.Parent()
	for found && p.Index() != 0 {
		path = append(path, p.Index())
		p, found = p.Parent()
	}
	return reverseInts(path)
}

// FindCommentsByDescriptor returns a `Comments` for the given descriptor.
// If not found, returns (nil, ErrCommentsNotFound).
func (s DescriptorSource) FindCommentsByDescriptor(d protoreflect.Descriptor) (Comments, error) {
	return s.FindCommentsByPath(getPath(d))
}

func reverseInts(path []int) []int {
	i, j := 0, len(path)-1
	for i < j {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

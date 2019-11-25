// Package descrule defines Protobuf descriptors
// and provides tools to extract them from a proto
// file.
package descrule

import (
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

type sourceInfo struct {
	leadingComments string
	file            *fileInfo
}

func (s *sourceInfo) LeadingComments() string {
	return s.leadingComments
}

func (s *sourceInfo) File() lint.FileInfo {
	return s.file
}

type fileInfo struct {
	path, comments string
}

func (f *fileInfo) Path() string {
	return f.path
}

func (f *fileInfo) Comments() string {
	return f.comments
}

func getSourceInfo(d desc.Descriptor) lint.SourceInfo {
	leadingComments := d.GetSourceInfo().GetLeadingComments()
	filePath := d.GetFile().GetName()
	fileComments := fileHeader(d.GetFile())
	return &sourceInfo{
		leadingComments: leadingComments,
		file: &fileInfo{
			path:     filePath,
			comments: fileComments,
		},
	}
}

// File is a wrapper around desc.FileDescriptor
type File struct {
	descriptor *desc.FileDescriptor
}

// NewFile creates a File wrapping desc.FileDescriptor.
func NewFile(f *desc.FileDescriptor) *File {
	return &File{descriptor: f}
}

// SourceInfo implements lint.ProtoRule.
func (f *File) SourceInfo() lint.SourceInfo {
	return getSourceInfo(f.descriptor)
}

// Message is a wrapper around desc.MessageDescriptor
type Message struct {
	descriptor *desc.MessageDescriptor
}

// NewMessage creates a Message wrapping a desc.MessageDescriptor
func NewMessage(m *desc.MessageDescriptor) *Message {
	return &Message{descriptor: m}
}

// SourceInfo implements lint.ProtoRule.
func (m *Message) SourceInfo() lint.SourceInfo {
	return getSourceInfo(m.descriptor)
}

// Enum is a wrapper around desc.EnumDescriptor
type Enum struct {
	descriptor *desc.EnumDescriptor
}

// NewEnum creates a Enum wrapping a desc.EnumDescriptor
func NewEnum(e *desc.EnumDescriptor) *Enum {
	return &Enum{descriptor: e}
}

// SourceInfo implements lint.ProtoRule.
func (e *Enum) SourceInfo() lint.SourceInfo {
	return getSourceInfo(e.descriptor)
}

// EnumValue is a wrapper around desc.EnumDescriptor
type EnumValue struct {
	descriptor *desc.EnumValueDescriptor
}

// NewEnumValue creates a EnumValue wrapping a desc.EnumValueDescriptor
func NewEnumValue(ev *desc.EnumValueDescriptor) *EnumValue {
	return &EnumValue{descriptor: ev}
}

// SourceInfo implements lint.ProtoRule.
func (ev *EnumValue) SourceInfo() lint.SourceInfo {
	return getSourceInfo(ev.descriptor)
}

// Service is a wrapper around desc.ServiceDescriptor
type Service struct {
	descriptor *desc.ServiceDescriptor
}

// NewService creates a Service wrapping a desc.ServiceDescriptor
func NewService(s *desc.ServiceDescriptor) *Service {
	return &Service{descriptor: s}
}

// SourceInfo implements lint.ProtoRule.
func (s *Service) SourceInfo() lint.SourceInfo {
	return getSourceInfo(s.descriptor)
}

// Method is a wrapper around desc.MethodDescriptor
type Method struct {
	descriptor *desc.MethodDescriptor
}

// NewMethod creates a Method wrapping a desc.MethodDescriptor
func NewMethod(m *desc.MethodDescriptor) *Method {
	return &Method{descriptor: m}
}

// SourceInfo implements lint.ProtoRule.
func (m *Method) SourceInfo() lint.SourceInfo {
	return getSourceInfo(m.descriptor)
}

// Field is a wrapper around desc.FieldDescriptor
type Field struct {
	descriptor *desc.FieldDescriptor
}

// NewField creates a Field wrapping a desc.FieldDescriptor
func NewField(f *desc.FieldDescriptor) *Field {
	return &Field{descriptor: f}
}

// SourceInfo implements lint.ProtoRule.
func (f *Field) SourceInfo() lint.SourceInfo {
	return getSourceInfo(f.descriptor)
}

// Descriptor is a wrapper around desc.Descriptor.
type Descriptor struct {
	descriptor desc.Descriptor
}

// NewDescriptor creates a Descriptor wrapping a desc.Descriptor
func NewDescriptor(d desc.Descriptor) *Descriptor {
	return &Descriptor{descriptor: d}
}

// SourceInfo implements lint.ProtoRule.
func (d *Descriptor) SourceInfo() lint.SourceInfo {
	return getSourceInfo(d.descriptor)
}

// AllDescriptors returns a slice with every descriptor in the file.
func AllDescriptors(f *desc.FileDescriptor) []lint.Descriptor {
	descriptors := []lint.Descriptor{&File{descriptor: f}}
	for _, d := range getAllEnums(f) {
		descriptors = append(descriptors, &Enum{descriptor: d})
	}
	for _, d := range getAllEnumValues(f) {
		descriptors = append(descriptors, &EnumValue{descriptor: d})
	}
	for _, d := range getAllFields(f) {
		descriptors = append(descriptors, &Field{descriptor: d})
	}
	for _, d := range getAllMethods(f) {
		descriptors = append(descriptors, &Method{descriptor: d})
	}
	for _, d := range getAllMessages(f) {
		descriptors = append(descriptors, &Message{descriptor: d})
	}
	for _, d := range f.GetServices() {
		descriptors = append(descriptors, &Service{descriptor: d})
	}
	return descriptors
}

// getAllMethods returns a slice with every method in the file.
func getAllMethods(f *desc.FileDescriptor) []*desc.MethodDescriptor {
	methods := []*desc.MethodDescriptor{}
	for _, s := range f.GetServices() {
		methods = append(methods, s.GetMethods()...)
	}
	return methods
}

// getAllMessages returns a slice with every message (not just top-level
// messages) in the file.
func getAllMessages(f *desc.FileDescriptor) (messages []*desc.MessageDescriptor) {
	messages = append(messages, f.GetMessageTypes()...)
	for _, message := range f.GetMessageTypes() {
		messages = append(messages, getAllNestedMessages(message)...)
	}
	return messages
}

// getAllNestedMessages returns a slice with the given message descriptor as well
// as all nested message descriptors, traversing to arbitrary depth.
func getAllNestedMessages(m *desc.MessageDescriptor) (messages []*desc.MessageDescriptor) {
	for _, nested := range m.GetNestedMessageTypes() {
		if !nested.IsMapEntry() { // Don't include the synthetic message type that represents an entry in a map field.
			messages = append(messages, nested)
		}
		messages = append(messages, getAllNestedMessages(nested)...)
	}
	return messages
}

// getAllEnums returns a slice with every enum (not just top-level enums)
// in the file.
func getAllEnums(f *desc.FileDescriptor) (enums []*desc.EnumDescriptor) {
	// Append all enums at the top level.
	enums = append(enums, f.GetEnumTypes()...)

	// Append all enums nested within messages.
	for _, m := range getAllMessages(f) {
		enums = append(enums, m.GetNestedEnumTypes()...)
	}

	return
}

// getAllEnumValues returns a slice with every enum value in the file.
func getAllEnumValues(f *desc.FileDescriptor) []*desc.EnumValueDescriptor {
	values := []*desc.EnumValueDescriptor{}
	for _, e := range getAllEnums(f) {
		values = append(values, e.GetValues()...)
	}
	return values
}

// getAllFields returns a slice with every field in the file.
func getAllFields(f *desc.FileDescriptor) []*desc.FieldDescriptor {
	fields := []*desc.FieldDescriptor{}
	for _, m := range getAllMessages(f) {
		fields = append(fields, m.GetFields()...)
	}
	return fields
}

// fileHeader attempts to get the comment at the top of the file, but it
// is on a best effort basis because protobuf is inconsistent.
//
// Taken from https://github.com/jhump/protoreflect/issues/215
func fileHeader(fd *desc.FileDescriptor) string {
	var firstLoc *dpb.SourceCodeInfo_Location
	var firstSpan int64

	// File level comments should only be comments identified on either
	// syntax (12), package (2), option (8), or import (3) statements.
	allowedPaths := map[int32]struct{}{2: {}, 3: {}, 8: {}, 12: {}}

	// Iterate over locations in the file descriptor looking for
	// what we think is a file-level comment.
	for _, curr := range fd.AsFileDescriptorProto().GetSourceCodeInfo().GetLocation() {
		// Skip locations that have no comments.
		if curr.LeadingComments == nil && len(curr.LeadingDetachedComments) == 0 {
			continue
		}
		// Skip locations that are not allowed because they should never be
		// mistaken for file-level comments.
		if _, ok := allowedPaths[curr.GetPath()[0]]; !ok {
			continue
		}
		currSpan := asPos(curr.Span)
		if firstLoc == nil || currSpan < firstSpan {
			firstLoc = curr
			firstSpan = currSpan
		}
	}
	if firstLoc == nil {
		return ""
	}
	if len(firstLoc.LeadingDetachedComments) > 0 {
		return strings.Join(firstLoc.LeadingDetachedComments, "\n")
	}
	return firstLoc.GetLeadingComments()
}

func asPos(span []int32) int64 {
	return (int64(span[0]) << 32) + int64(span[1])
}

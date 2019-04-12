package lint

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/google/go-cmp/cmp"
)

//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_source.protoset --proto_path=testdata testdata/test_source.proto
//go:generate protoc --include_source_info --descriptor_set_out=testdata/test_rule_disable.protoset --proto_path=testdata testdata/test_rule_disable.proto

type testDescriptorVisiting struct {
	visit func(d protoreflect.Descriptor)
}

func (v testDescriptorVisiting) VisitDescriptor(d protoreflect.Descriptor) {
	v.visit(d)
}

func TestDescriptorLocation(t *testing.T) {
	fileDesc, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

	tests := []struct {
		descriptor protoreflect.Descriptor
		want       Location
	}{
		{
			descriptor: fileDesc.Messages().Get(0), // A top level message.
			want: Location{
				Start: Position{
					Line: 7, Column: 0,
				},
				End: Position{
					Line: 59, Column: 1,
				},
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Messages().Get(0), // A nested message.
			want: Location{
				Start: Position{
					Line: 9, Column: 2,
				},
				End: Position{
					Line: 36, Column: 3,
				},
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0),
			want: Location{
				Start: Position{
					Line: 44, Column: 2,
				},
				End: Position{
					Line: 49, Column: 3,
				},
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0).Values().Get(1),
			want: Location{
				Start: Position{
					Line: 48, Column: 4,
				},
				End: Position{
					Line: 48, Column: 12,
				},
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Fields().Get(1),
			want: Location{
				Start: Position{
					Line: 41, Column: 2,
				},
				End: Position{
					Line: 41, Column: 41,
				},
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Oneofs().Get(0),
			want: Location{
				Start: Position{
					Line: 53, Column: 2,
				},
				End: Position{
					Line: 58, Column: 3,
				},
			},
		},
		{
			descriptor: fileDesc.Services().Get(0),
			want: Location{
				Start: Position{
					Line: 72, Column: 0,
				},
				End: Position{
					Line: 77, Column: 1,
				},
			},
		},
		{
			descriptor: fileDesc.Services().Get(0).Methods().Get(0),
			want: Location{
				Start: Position{
					Line: 74, Column: 2,
				},
				End: Position{
					Line: 74, Column: 44,
				},
			},
		},
	}

	for _, test := range tests {
		got, err := descSource.DescriptorLocation(test.descriptor)
		if err != nil {
			t.Errorf("DescriptorLocation(%s) error: %s", test.descriptor.FullName(), err)
		}
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("DescriptorLocation(%s) mismatch (-want +got):\n%s", test.descriptor.FullName(), diff)
		}
	}
}

func TestDescriptorComments(t *testing.T) {
	fileDesc, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

	tests := []struct {
		descriptor protoreflect.Descriptor
		want       Comments
	}{
		{
			descriptor: fileDesc.Messages().Get(0),
			want: Comments{
				LeadingComments: " Outer\n",
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Enums().Get(0),
			want: Comments{
				LeadingComments: " NestedEnum\n",
			},
		},
		{
			descriptor: fileDesc.Messages().Get(0).Fields().Get(1),
			want: Comments{
				LeadingComments: " outer_middle_field\n",
			},
		},
	}

	for _, test := range tests {
		got, err := descSource.DescriptorComments(test.descriptor)
		if err != nil {
			t.Errorf("DescriptorComments(%s) error: %s", test.descriptor.FullName(), err)
		}
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("DescriptorComments(%s) mismatch (-want +got):\n%s", test.descriptor.FullName(), diff)
		}
	}
}

func TestSyntaxLocation(t *testing.T) {
	_, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

	want := Location{
		Start: Position{Line: 2, Column: 0},
		End:   Position{Line: 2, Column: 18},
	}
	got, err := descSource.SyntaxLocation()
	if err != nil {
		t.Errorf("SyntaxLocation() error: %s", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("SyntaxLocation() mismatch (-want +got):\n%s", diff)
	}
}

func TestSyntaxComments(t *testing.T) {
	_, proto := readProtoFile("test_source.protoset")
	descSource, err := newDescriptorSource(proto)
	if err != nil {
		t.Errorf("newDescriptorSource: %v", err)
	}

	want := Comments{
		LeadingDetachedComments: []string{" DO NOT EDIT -- This is for `source_test.go`.\n"},
	}
	got, err := descSource.SyntaxComments()
	if err != nil {
		t.Errorf("SyntaxComments() error: %s", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("SyntaxComments() mismatch (-want +got):\n%s", diff)
	}
}

func readProtoFile(fileName string) (protoreflect.FileDescriptor, *descriptorpb.FileDescriptorProto) {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	proto := protoset.GetFile()[0]
	f, err := protodesc.NewFile(proto, nil)
	if err != nil {
		log.Fatalf("protodesc.NewFile() error: %v", err)
	}
	return f, proto
}

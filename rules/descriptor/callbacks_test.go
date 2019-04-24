package descriptor

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	"github.com/jgeewax/api-linter/lint"
)

func TestCallbacks_Apply(t *testing.T) {
	f := readProtoFile("test.protoset")
	descriptors := map[string]protoreflect.Descriptor{
		"enum":      f.Enums().Get(0),
		"enumvalue": f.Enums().Get(0).Values().Get(0),
		"field":     f.Messages().Get(0).Fields().Get(0),
		"message":   f.Messages().Get(0),
		"method":    f.Services().Get(0).Methods().Get(0),
		"oneof":     f.Messages().Get(0).Oneofs().Get(0),
		"service":   f.Services().Get(0),
	}

	tests := []struct {
		callbacks Callbacks
		results   []lint.Problem
	}{
		{Callbacks{}, []lint.Problem{}},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "enum"}}, nil
				},
			},
			[]lint.Problem{{Message: "enum"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "enumvalue"}}, nil
				},
			},
			[]lint.Problem{{Message: "enumvalue"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "field"}}, nil
				},
			},
			[]lint.Problem{{Message: "field"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "message"}}, nil
				},
			},
			[]lint.Problem{{Message: "message"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "method"}}, nil
				},
			},
			[]lint.Problem{{Message: "method"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "oneof"}}, nil
				},
			},
			[]lint.Problem{{Message: "oneof"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "service"}}, nil
				},
			},
			[]lint.Problem{{Message: "service"}},
		},
	}

	for _, test := range tests {
		results := []lint.Problem{}
		for _, d := range descriptors {
			problems, err := test.callbacks.Apply(d, lint.DescriptorSource{})
			if err != nil {
				t.Errorf("Callbacks.Apply returns unexpected error: %v", err)
			}
			results = append(results, problems...)
		}
		if got, want := results, test.results; !reflect.DeepEqual(got, want) {
			t.Errorf("Callbacks.Apply returns problems '%s', but want '%s'", got, want)
		}
	}
}

func TestCallbacks_Apply_DescriptorCallback(t *testing.T) {
	f := readProtoFile("test.protoset")
	descriptors := map[string]protoreflect.Descriptor{
		"enum":      f.Enums().Get(0),
		"enumvalue": f.Enums().Get(0).Values().Get(0),
	}

	tests := []struct {
		callbacks  Callbacks
		numProblem int
	}{
		{Callbacks{}, 0},
		{
			Callbacks{
				DescriptorCallback: func(d protoreflect.Descriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "desc"}}, nil
				},
			},
			2,
		},
	}

	for _, test := range tests {
		all := []lint.Problem{}
		for _, d := range descriptors {
			problems, err := test.callbacks.Apply(d, lint.DescriptorSource{})
			if err != nil {
				t.Errorf("Callbacks.Apply returns unexpected error: %v", err)
			}
			all = append(all, problems...)
		}
		if got, want := len(all), test.numProblem; got != want {
			t.Errorf("Callbacks.Apply returns %d problems, but want %d", got, want)
		}
	}
}

func readProtoFile(fileName string) protoreflect.FileDescriptor {
	path := filepath.Join("testdata", fileName)
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %v", path, err)
	}
	protoset := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bs, protoset); err != nil {
		log.Fatalf("Unable to parse %T from %s: %v", protoset, path, err)
	}
	f, err := protodesc.NewFile(protoset.GetFile()[0], nil)
	if err != nil {
		log.Fatalf("protodesc.NewFile() error: %v", err)
	}
	return f

}

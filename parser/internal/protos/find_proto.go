package protos

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	protoV1 "github.com/golang/protobuf/proto"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FindByName returned a FileDescriptorProto if the given name is found
// in the registered proto files.
func FindByName(filename string) (*descriptorpb.FileDescriptorProto, error) {
	d := protoV1.FileDescriptor(filename)
	if d == nil {
		return nil, fmt.Errorf("Proto file %q is not found", filename)
	}
	b := bytes.NewBuffer(d)
	g, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	c, err := ioutil.ReadAll(g)
	if err != nil {
		return nil, err
	}
	f := &descriptorpb.FileDescriptorProto{}
	err = protoV2.Unmarshal(c, f)
	return f, err
}

// FindAllByNames returns a list of FileDescriptorProto if all the names are found
// in the registered proto files.
func FindAllByNames(filename ...string) ([]*descriptorpb.FileDescriptorProto, error) {
	var files []*descriptorpb.FileDescriptorProto
	for _, name := range filename {
		f, err := FindByName(name)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

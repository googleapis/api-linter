package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

type protocParser struct {
	protoc     string
	importPath string
}

func (p *protocParser) ParseProto(filenames ...string) ([]*descriptorpb.FileDescriptorProto, error) {
	return compileProto(p.protoc, p.importPath, filenames...)
}

func compileProto(protoc, importPath string, filenames ...string) ([]*descriptorpb.FileDescriptorProto, error) {
	outfile, err := ioutil.TempFile("", "desc.pb")
	if err != nil {
		return nil, err
	}
	defer os.Remove(outfile.Name())

	cmd := exec.Command(protoc, "--descriptor_set_out="+outfile.Name())
	cmd.Args = append(cmd.Args, "--include_source_info")
	cmd.Args = append(cmd.Args, "-I="+importPath)
	cmd.Args = append(cmd.Args, filenames...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("Running %s: %s, %v", strings.Join(cmd.Args, " "), string(out), err)
	}

	b, err := ioutil.ReadFile(outfile.Name())
	if err != nil {
		return nil, err
	}
	v := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v.GetFile(), nil
}

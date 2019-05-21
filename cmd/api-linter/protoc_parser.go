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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
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
	cmd.Args = append(cmd.Args, "--proto_path="+importPath)
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

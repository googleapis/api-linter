package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/v2/proto"
	"github.com/golang/protobuf/v2/reflect/protodesc"
	"github.com/golang/protobuf/v2/reflect/protoreflect"
	"github.com/golang/protobuf/v2/reflect/protoregistry"
	pluginpb "github.com/golang/protobuf/v2/types/plugin"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/rules"
)

func main() {
	rules := rules.Rules()
	log.Printf("Number of rules: %d\n", len(rules.AllRules()))
	linter := newLinter(*rules)
	linter.run()
}

type linter struct {
	reader io.Reader
	writer io.Writer
	rules  lint.Rules
}

func (l linter) run() {
	request := readLintRequest(l.reader)
	response, err := lint.Run(l.rules, request)
	if err != nil {
		log.Fatalf("Error when running lint: %v", err)
	}
	log.Printf("Total number of API-Linter findings: %d", len(response.Problems))
	for _, problem := range response.Problems {
		log.Printf("Finding: %s, suggestion: %s", problem.Message, problem.Suggestion)
	}
}

func newLinter(rules lint.Rules) *linter {
	return &linter{
		rules:  rules,
		reader: os.Stdin,
		writer: os.Stdout,
	}
}

func readLintRequest(r io.Reader) lint.Request {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("Error when reading CodeGeneratorRequest: %v", err)
	}

	var codeGenRequest pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(in, &codeGenRequest); err != nil {
		log.Fatalf("Error when unmarshaling CodeGeneratorRequest: %v", err)
	}

	if len(codeGenRequest.GetProtoFile()) == 0 {
		log.Fatalf("Error: zero proto files in CodeGeneratorRequest")
	}

	fd := codeGenRequest.GetProtoFile()[0]
	f, err := protodesc.NewFile(fd, protoregistry.NewFiles())
	if err != nil {
		log.Fatalf("Error when converting proto to descriptor: %v", err)
	}

	ctx := context.Background()
	source, err := lint.NewDescriptorSource(fd)
	if err != nil {
		return lintRequest{
			protoFile: f,
			context:   lint.NewContext(ctx),
		}
	}

	return lintRequest{
		protoFile: f,
		context:   lint.NewContextWithDescriptorSource(ctx, source),
	}
}

type lintRequest struct {
	protoFile protoreflect.FileDescriptor
	context   lint.Context
}

func (r lintRequest) ProtoFile() protoreflect.FileDescriptor {
	return r.protoFile
}

func (r lintRequest) Context() lint.Context {
	return r.context
}

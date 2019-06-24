// Package protoc provides a parser wrapping protobuf compiler to
// parse protobuf files.
package protoc

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/parser/internal/protos"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Option defines an option applied to the `protoc` parser.
type Option func(p *Parser)

var defaultCommandPath = func(p *Parser) {
	p.command = "protoc"
}

// Parser is a wrapper of protobuf compiler that parses protobuf files.
// By default, it includes source_code_info, and automatically
// resolve common protos.
type Parser struct {
	command             string
	protoPaths          []string
	descriptorSetIn     []*descriptorpb.FileDescriptorProto
	includeImports      bool
	excludeCommonProtos bool
}

// New creates a new Parser with option functions.
func New(options ...Option) *Parser {
	p := &Parser{}
	defaultCommandPath(p)
	for _, o := range options {
		o(p)
	}
	return p
}

// IncludeImports sets `--include_imports` flag to `protoc`.
func IncludeImports() Option {
	return func(p *Parser) {
		p.includeImports = true
	}
}

// AddProtoPath adds `--proto_path` flag to `protoc`
func AddProtoPath(path ...string) Option {
	return func(p *Parser) {
		p.protoPaths = append(p.protoPaths, path...)
	}
}

// AddDescriptorSetIn adds `--descriptor_set_in` to `protoc`.
func AddDescriptorSetIn(files ...*descriptorpb.FileDescriptorProto) Option {
	return func(p *Parser) {
		p.descriptorSetIn = append(p.descriptorSetIn, files...)
	}
}

// ExcludeCommonProtos excludes pre-loaded, common protobuf files.
func ExcludeCommonProtos() Option {
	return func(p *Parser) {
		p.excludeCommonProtos = true
	}
}

// Command sets the `protoc` command line path.
func Command(path string) Option {
	return func(p *Parser) {
		p.command = path
	}
}

// Parse parses a list of protobuf file with the given Options.
func (p *Parser) Parse(files ...string) (*descriptorpb.FileDescriptorSet, error) {
	return p.parse(files...)
}

func (p *Parser) parse(files ...string) (*descriptorpb.FileDescriptorSet, error) {
	args := []string{
		"--include_source_info",
		"--proto_path=.",
	}
	if p.includeImports {
		args = append(args, "--include_imports")
	}

	for _, path := range p.protoPaths {
		args = append(args, "--proto_path="+path)
	}

	fileDescSet := &descriptorpb.FileDescriptorSet{}
	fileDescSet.File = append(fileDescSet.File, p.descriptorSetIn...)
	if !p.excludeCommonProtos {
		imports, err := extractImportsFromFiles(files...)
		if err != nil {
			return nil, err
		}
		for _, name := range imports {
			if f, err := protos.FindByName(name); err == nil {
				fileDescSet.File = append(fileDescSet.File, f)
			}
		}
	}
	if len(fileDescSet.GetFile()) > 0 {
		descSetInFile, err := ioutil.TempFile("", "protoc_compiler")
		if err != nil {
			return nil, err
		}
		defer os.Remove(descSetInFile.Name())
		writeFileDescriptorSet(descSetInFile.Name(), fileDescSet)
		args = append(args, "--descriptor_set_in="+descSetInFile.Name())
	}

	tempOutFile, err := ioutil.TempFile("", "protoc_compiler")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempOutFile.Name())
	args = append(args, "--descriptor_set_out="+tempOutFile.Name())

	args = append(args, files...)

	cmd := exec.Command(p.command, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("Running %s: %s, %v", strings.Join(cmd.Args, " "), string(out), err)
	}

	b, err := ioutil.ReadFile(tempOutFile.Name())
	if err != nil {
		return nil, err
	}
	out := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(b, out); err != nil {
		return nil, err
	}
	return out, err
}

func extractImportsFromFiles(files ...string) ([]string, error) {
	s := make(map[string]struct{})
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		names := extractImports(string(b))
		for _, name := range names {
			s[name] = struct{}{}
		}
	}
	var names []string
	for key := range s {
		names = append(names, key)
	}
	return names, nil
}

var importPattern = regexp.MustCompile(`(?m)^\s*import\s+"([a-zA-Z0-9/\.]+)"\s*;`)

func extractImports(s string) []string {
	var imports []string
	findings := importPattern.FindAllStringSubmatch(s, -1)
	for _, f := range findings {
		imports = append(imports, f[1])
	}
	return imports
}

func writeFileDescriptorSet(outfile string, s *descriptorpb.FileDescriptorSet) error {
	b, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outfile, b, 0666)
}

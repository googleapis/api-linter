package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
	pluginpb "github.com/golang/protobuf/v2/types/plugin"
)

// Options are optional parameters to NewPlugin.
type Options struct {
	// If ParamFunc is non-nil, it will be called with each unknown
	// generator parameter.
	//
	// Plugins for protoc can accept parameters from the command line,
	// passed in the --<plugin_name>_out protoc, separated from the output
	// directory with a colon; e.g.,
	//
	//   --api_linter_out=<param1>=<value1>,<param2>=<value2>:<output_directory>
	//
	// Parameters passed in this fashion as a comma-separated list of
	// key=value pairs will be passed to the ParamFunc.
	//
	// The (flag.FlagSet).Set method matches this function signature,
	// so parameters can be converted into flags as in the following:
	//
	//   var flags flag.FlagSet
	//   value := flags.Bool("param", false, "")
	//   opts := &protogen.Options{
	//     ParamFunc: flags.Set,
	//   }
	//   protogen.Run(opts, func(p *protogen.Plugin) error {
	//     if *value { ... }
	//   })
	ParamFunc func(name, value string) error
}

// A Plugin is a protoc plugin that generates or checks files.
type Plugin struct {
	filesToGenerate []*descriptorpb.FileDescriptorProto
	genFiles        []*GeneratedFile
	err             error
}

// NewPlugin returns a new Plugin.
func NewPlugin(req *pluginpb.CodeGeneratorRequest, opts *Options) (*Plugin, error) {
	gen := &Plugin{}

	gen.extractFilesToGenerate(req)
	if err := gen.extractParameters(req, opts); err != nil {
		return nil, err
	}

	return gen, nil
}

func (gen *Plugin) extractFilesToGenerate(req *pluginpb.CodeGeneratorRequest) {
	filesToGenerate := make(map[string]bool)
	for _, f := range req.GetFileToGenerate() {
		filesToGenerate[f] = true
	}
	for _, f := range req.GetProtoFile() {
		if filesToGenerate[f.GetName()] {
			gen.filesToGenerate = append(gen.filesToGenerate, f)
		}
	}
}

func (gen *Plugin) extractParameters(req *pluginpb.CodeGeneratorRequest, opts *Options) error {
	for _, param := range strings.Split(req.GetParameter(), ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		if param != "" {
			if err := opts.ParamFunc(param, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// Error records an error during running this plugin. The plugin will
// report the error back to protoc and will not produce output.
func (gen *Plugin) Error(err error) {
	if gen.err == nil && err != nil {
		gen.err = err
	}
}

// Files returns files to be generated or checked.
func (gen *Plugin) Files() []*descriptorpb.FileDescriptorProto {
	return gen.filesToGenerate
}

// Response returns the generator output.
func (gen *Plugin) Response() *pluginpb.CodeGeneratorResponse {
	resp := &pluginpb.CodeGeneratorResponse{}
	if gen.err != nil {
		resp.Error = proto.String(gen.err.Error())
		return resp
	}

	for _, f := range gen.genFiles {
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(f.filename),
			Content: proto.String(string(f.Content())),
		})
	}

	return resp
}

// NewGeneratedFile creates a new generated file with the given filename.
func (gen *Plugin) NewGeneratedFile(filename string) *GeneratedFile {
	f := &GeneratedFile{filename: filename}
	gen.genFiles = append(gen.genFiles, f)
	return f
}

// Run executes a function as a protoc plugin.
//
// It reads a CodeGeneratorRequest message from os.Stdin, invokes the plugin
// function, and writes a CodeGeneratorResponse message to os.Stdout.
//
// If a failure occurs while reading or writing, Run prints an error to
// os.Stderr and calls os.Exit(1).
func Run(opts *Options, f func(*Plugin) error) {
	if err := run(opts, f); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func run(opts *Options, f func(*Plugin) error) error {
	req, err := readCodeGeneratorRequest(os.Stdin)
	if err != nil {
		return err
	}
	gen, err := NewPlugin(req, opts)
	if err != nil {
		return err
	}

	if err := f(gen); err != nil {
		// Errors from the plugin function are reported by setting
		// the error field in the CodeGeneratorResponse.
		//
		// In contrast, errors that indicate a problem in protoc
		// itself (un-parsable input, I/O errors, etc.) are reported
		// to stderr
		gen.Error(err)
	}

	resp := gen.Response()
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}

func readCodeGeneratorRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return nil, err
	}
	return req, nil
}

// A GeneratedFile is a generated file.
type GeneratedFile struct {
	filename string
	buf      bytes.Buffer
}

// Write writes contents to the generated file.
func (f *GeneratedFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

// Content returns the contents of the generated file.
func (f *GeneratedFile) Content() []byte {
	return f.buf.Bytes()
}

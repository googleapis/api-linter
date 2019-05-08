// Package protogen provides support for writing plugins for protoc,
// the Protocol Buffers Compiler. The protoc plugins read a
// CodeGeneratorRequest protocol buffer from standard input and
// write a CodeGeneratorResponse to standard output.
package protogen

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
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
	ParamFunc func(name, value string) error
}

// A Plugin is a protoc plugin that reads a CodeGeneratorRequest protocol buffer
// from standard input and writes a CodeGeneratorResponse protocol buffer to
// standard output.
type Plugin struct {
	req             *pluginpb.CodeGeneratorRequest
	filesToGenerate []*descriptorpb.FileDescriptorProto
	genFiles        []*GeneratedFile
	err             error
}

// NewPlugin returns a new Plugin.
func NewPlugin(req *pluginpb.CodeGeneratorRequest, opts *Options) (*Plugin, error) {
	gen := &Plugin{req: req}

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

// Files returns files to be generated, which don't include the imported files.
func (gen *Plugin) Files() []*descriptorpb.FileDescriptorProto {
	return gen.filesToGenerate
}

// Request returns a CodeGeneratorRequest.
func (gen *Plugin) Request() *pluginpb.CodeGeneratorRequest {
	return gen.req
}

// Response returns a CodeGeneratorResponse.
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

// Generator provides Options for creating a Plugin and generates
// files by using the information provided by the Plugin.
type Generator interface {
	// Options returns a Options for NewPlugin.
	Options() *Options
	// Generate consumes a Plugin and generates(checks) files.
	Generate(*Plugin) error
}

// Run executes the given Generator.
//
// It reads a CodeGeneratorRequest message from os.Stdin, invokes the generator,
// and writes a CodeGeneratorResponse message to os.Stdout.
//
// If a failure occurs while reading or writing, Run prints an error to
// os.Stderr and calls os.Exit(1).
func Run(g Generator) error {
	return run(g)
}

func run(g Generator) error {
	req, err := readCodeGeneratorRequest(os.Stdin)
	if err != nil {
		return err
	}
	gen, err := NewPlugin(req, g.Options())
	if err != nil {
		return err
	}

	if err := g.Generate(gen); err != nil {
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

// GeneratedFile stores the filename and content
// of a file that will be generated by protoc.
type GeneratedFile struct {
	filename string
	buf      bytes.Buffer
}

// Write writes bytes p.
func (f *GeneratedFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

// Content returns the contents of the generated file.
func (f *GeneratedFile) Content() []byte {
	return f.buf.Bytes()
}

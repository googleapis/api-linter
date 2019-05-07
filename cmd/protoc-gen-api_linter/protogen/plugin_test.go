package protogen_test

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	pluginpb "github.com/golang/protobuf/v2/types/plugin"
	"github.com/googleapis/api-linter/cmd/protoc-gen-api_linter/protogen"
)

func TestPluginParameters(t *testing.T) {
	var flags flag.FlagSet
	value := flags.Int("integer", 0, "")
	opts := &protogen.Options{
		ParamFunc: flags.Set,
	}
	const params = "integer=2"
	_, err := protogen.NewPlugin(&pluginpb.CodeGeneratorRequest{
		Parameter: proto.String(params),
	}, opts)
	if err != nil {
		t.Errorf("NewPlugin(generator parameters %q): %v", params, err)
	}
	if *value != 2 {
		t.Errorf("NewPlugin(generator parameters %q): integer=%v, want 2", params, *value)
	}
}

func TestPluginParameterErrors(t *testing.T) {
	for _, parameter := range []string{
		"unknown=1",
		"boolean=error",
	} {
		var flags flag.FlagSet
		flags.Bool("boolean", false, "")
		opts := &protogen.Options{
			ParamFunc: flags.Set,
		}
		_, err := protogen.NewPlugin(&pluginpb.CodeGeneratorRequest{
			Parameter: proto.String(parameter),
		}, opts)
		if err == nil {
			t.Errorf("New(generator parameters %q): want error, got nil", parameter)
		}
	}
}

func TestFiles(t *testing.T) {
	gen, err := protogen.NewPlugin(makeRequest(t, "testdata/test.proto"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(gen.Files()) != 1 {
		t.Errorf("Files() returns %d files, but want only 1", len(gen.Files()))
	}
	if got, want := gen.Files()[0].GetName(), "test.proto"; got != want {
		t.Errorf("Files() returns file %q, but want %q", got, want)
	}
}

func TestResponseWithError(t *testing.T) {
	tests := []struct {
		err     error
		respErr string
	}{
		{nil, ""},
		{errors.New("testing"), "testing"},
	}
	for _, test := range tests {
		gen, err := protogen.NewPlugin(&pluginpb.CodeGeneratorRequest{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		gen.Error(test.err)
		resp := gen.Response()
		if got, want := resp.GetError(), test.respErr; got != want {
			t.Errorf("Response got %q error, but want %q", got, want)
		}
	}
}

func TestResponseWithGeneratedFiles(t *testing.T) {
	type file struct {
		name, content string
	}
	tests := []struct {
		inputs, wants []file
	}{
		{
			inputs: []file{},
			wants:  []file{},
		},
		{
			inputs: []file{
				{
					name:    "test_name",
					content: "test_content",
				},
			},
			wants: []file{
				{
					name:    "test_name",
					content: "test_content",
				},
			},
		},
		{
			inputs: []file{
				{
					name:    "test_name",
					content: "test_content",
				},
				{
					name:    "test_name_2",
					content: "test_content_2",
				},
			},
			wants: []file{
				{
					name:    "test_name",
					content: "test_content",
				},
				{
					name:    "test_name_2",
					content: "test_content_2",
				},
			},
		},
	}
	for _, test := range tests {
		gen, err := protogen.NewPlugin(&pluginpb.CodeGeneratorRequest{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range test.inputs {
			gen.NewGeneratedFile(f.name).Write([]byte(f.content))
		}
		resp := gen.Response()
		if resp.GetError() != "" {
			t.Errorf("Unexpected error: %q", resp.GetError())
		}
		results := []file{}
		for _, f := range resp.GetFile() {
			results = append(results, file{f.GetName(), f.GetContent()})
		}
		if !reflect.DeepEqual(results, test.wants) {
			t.Errorf("Generated files %v, but want %v", results, test.wants)
		}
	}
}

// makeRequest returns a CodeGeneratorRequest for the given protoc inputs.
//
// It does this by running protoc with the current binary as the protoc-gen-go
// plugin. This "plugin" produces a single file, named 'request', which contains
// the code generator request.
func makeRequest(t *testing.T, args ...string) *pluginpb.CodeGeneratorRequest {
	workdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	cmd := exec.Command("protoc", "--plugin=protoc-gen-protogen="+os.Args[0])
	cmd.Args = append(cmd.Args, "--protogen_out="+workdir, "-Itestdata")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(os.Environ(), "RUN_AS_PROTOC_PLUGIN=1")
	out, err := cmd.CombinedOutput()
	if len(out) > 0 || err != nil {
		t.Log("RUNNING: ", strings.Join(cmd.Args, " "))
	}
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Fatalf("protoc: %v", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(workdir, "request"))
	if err != nil {
		t.Fatal(err)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.UnmarshalText(string(b), req); err != nil {
		t.Fatal(err)
	}
	return req
}

func init() {
	if os.Getenv("RUN_AS_PROTOC_PLUGIN") != "" {
		protogen.Run(&testGenerator{})
		os.Exit(0)
	}
}

type testGenerator struct{}

func (t *testGenerator) Options() *protogen.Options {
	return nil
}

func (t *testGenerator) Generate(p *protogen.Plugin) error {
	g := p.NewGeneratedFile("request")
	return proto.MarshalText(g, p.Request())
}

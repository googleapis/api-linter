package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"gopkg.in/yaml.v2"
)

type cli struct {
	configPath    string
	formatType    string
	outputPath    string
	protoImports  []string
	protoFiles    []string
	protoDescPath string
}

func newCli(args []string) *cli {
	// Define flags.
	var fmtFlag string
	var cfgFlag string
	var outFlag string
	var protoImportFlag stringSlice
	var protoDescFlag string

	// Register flags.
	fs := flag.NewFlagSet("api-linter", flag.ExitOnError)
	fs.Var(&protoImportFlag, "proto_path", "The folder to search for proto imports.")
	fs.StringVar(&cfgFlag, "config", "", "The linter config file.")
	fs.StringVar(&fmtFlag, "out_format", "", "The format of the linting results.")
	fs.StringVar(&outFlag, "out_path", "", "The output file path.")
	fs.StringVar(&protoDescFlag, "proto_desc", "", "The file descriptor set for proto imports")

	// Parse flags.
	fs.Parse(args)

	c := &cli{
		configPath:    cfgFlag,
		formatType:    fmtFlag,
		outputPath:    outFlag,
		protoImports:  []string{"."},
		protoFiles:    fs.Args(),
		protoDescPath: protoDescFlag,
	}

	c.protoImports = append(c.protoImports, protoImportFlag...)

	return c
}

func (c *cli) lint(rules lint.RuleRegistry, configs lint.Configs) error {
	// Check if there are files to lint; if not, abort.
	if len(c.protoFiles) == 0 {
		return fmt.Errorf("no files to lint")
	}
	// Prepare proto import lookup.
	var lookupImport func(string) (*desc.FileDescriptor, error)
	if c.protoDescPath != "" {
		fs, err := loadFileDescriptors(c.protoDescPath)
		if err != nil {
			return err
		}
		lookupImport = func(name string) (*desc.FileDescriptor, error) {
			if f, found := fs[name]; found {
				return f, nil
			}
			return nil, fmt.Errorf("%q is not found", name)
		}
	}
	// Parse files into protoreflect file descriptors.
	p := protoparse.Parser{
		ImportPaths:           c.protoImports,
		IncludeSourceCodeInfo: true,
		LookupImport:          lookupImport,
	}
	fd, err := p.ParseFiles(c.protoFiles...)
	if err != nil {
		return err
	}

	if c.configPath != "" {
		config, err := lint.ReadConfigsFromFile(c.configPath)
		if err != nil {
			return err
		}
		configs = append(configs, config...)
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rules, configs)
	results, err := l.LintProtos(fd...)
	if err != nil {
		return err
	}

	// Determine the output to write the results.
	// Stdout is the default output.
	w := os.Stdout
	if c.outputPath != "" {
		var err error
		w, err = os.Create(c.outputPath)
		if err != nil {
			return err
		}
		defer w.Close()
	}

	// Determine the format to print the results.
	// YAML format is the default.
	marshal := yaml.Marshal
	switch c.formatType {
	case "json":
		marshal = json.Marshal
	case "summary":
		marshal = func(i interface{}) ([]byte, error) {
			return emitSummary(i.([]lint.Response))
		}
	}

	// Print the results.
	b, err := marshal(results)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

type stringSlice []string

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (p *stringSlice) String() string {
	return fmt.Sprint(*p)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (p *stringSlice) Set(value string) error {
	*p = append(*p, strings.Split(value, ",")...)
	return nil
}

func loadFileDescriptors(filePath string) (map[string]*desc.FileDescriptor, error) {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorsFromSet(fs)
}

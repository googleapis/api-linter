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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type cli struct {
	ConfigPath       string
	FormatType       string
	OutputPath       string
	ProtoImportPaths []string
	ProtoFiles       []string
	ProtoDescPath    string
	EnabledRules     []string
	DisabledRules    []string
}

func newCli(args []string) *cli {
	// Define flag variables.
	var cfgFlag string
	var fmtFlag string
	var outFlag string
	var protoImportFlag []string
	var protoDescFlag string
	var ruleEnableFlag []string
	var ruleDisableFlag []string

	// Register flag variables.
	fs := pflag.NewFlagSet("api-linter", pflag.ExitOnError)
	fs.StringVar(&cfgFlag, "config", "", "The linter config file.")
	fs.StringVar(&fmtFlag, "output-format", "", "The format of the linting results.\nSupported formats include \"yaml\", \"json\" and \"summary\" table.\nYAML is the default.")
	fs.StringVarP(&outFlag, "output-path", "o", "", "The output file path.\nIf not given, the linting results will be printed out to STDOUT.")
	fs.StringArrayVarP(&protoImportFlag, "proto-path", "I", nil, "The folder for searching proto imports.\nMay be specified multiple times; directories will be searched in order.\nThe current working directory is always used.")
	fs.StringVar(&protoDescFlag, "proto-descriptor-set", "", "A delimited (':') list of files each containing a FileDescriptorSet for searching proto imports.")
	fs.StringArrayVar(&ruleEnableFlag, "enable-rule", nil, "Enable a rule with the given name.\nMay be specified multiple times.")
	fs.StringArrayVar(&ruleDisableFlag, "disable-rule", nil, "Disable a rule with the given name.\nMay be specified multiple times.")

	// Parse flags.
	fs.Parse(args)

	return &cli{
		ConfigPath:       cfgFlag,
		FormatType:       fmtFlag,
		OutputPath:       outFlag,
		ProtoImportPaths: append(protoImportFlag, "."),
		ProtoDescPath:    protoDescFlag,
		EnabledRules:     ruleEnableFlag,
		DisabledRules:    ruleDisableFlag,
		ProtoFiles:       fs.Args(),
	}
}

func (c *cli) lint(rules lint.RuleRegistry, configs lint.Configs) error {
	// Pre-check if there are files to lint.
	if len(c.ProtoFiles) == 0 {
		return fmt.Errorf("no file to lint")
	}
	// Read linter config and append it to the default.
	if c.ConfigPath != "" {
		config, err := lint.ReadConfigsFromFile(c.ConfigPath)
		if err != nil {
			return err
		}
		configs = append(configs, config...)
	}
	// Add configs for the enabled rules.
	for _, ruleName := range c.EnabledRules {
		configs = append(configs, lint.Config{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				ruleName: {}, // default is enabled.
			},
		})
	}
	// Add configs for the disabled rules.
	for _, ruleName := range c.DisabledRules {
		configs = append(configs, lint.Config{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				ruleName: {Disabled: true},
			},
		})
	}
	// Prepare proto import lookup.
	var lookupImport func(string) (*desc.FileDescriptor, error)
	if c.ProtoDescPath != "" {
		fs, err := loadFileDescriptors(strings.Split(c.ProtoDescPath, ":")...)
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
	// Parse proto files into `protoreflect` file descriptors.
	p := protoparse.Parser{
		ImportPaths:           c.ProtoImportPaths,
		IncludeSourceCodeInfo: true,
		LookupImport:          lookupImport,
	}
	// Resolve file absolute paths to relative ones.
	protoFiles, err := protoparse.ResolveFilenames(c.ProtoImportPaths, c.ProtoFiles...)
	if err != nil {
		return err
	}
	fd, err := p.ParseFiles(protoFiles...)
	if err != nil {
		return err
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rules, configs)
	results, err := l.LintProtos(fd...)
	if err != nil {
		return err
	}

	// Determine the output for writing the results.
	// Stdout is the default output.
	w := os.Stdout
	if c.OutputPath != "" {
		var err error
		w, err = os.Create(c.OutputPath)
		if err != nil {
			return err
		}
		defer w.Close()
	}

	// Determine the format for printing the results.
	// YAML format is the default.
	marshal := yaml.Marshal
	switch c.FormatType {
	case "json":
		marshal = json.Marshal
	case "summary":
		marshal = func(i interface{}) ([]byte, error) {
			return printSummaryTable(i.([]lint.Response))
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

func loadFileDescriptors(filePaths ...string) (map[string]*desc.FileDescriptor, error) {
	fds := []*dpb.FileDescriptorProto{}
	for _, filePath := range filePaths {
		fs, err := readFileDescriptorSet(filePath)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fs.GetFile()...)
	}
	return desc.CreateFileDescriptors(fds)
}

func readFileDescriptorSet(filePath string) (*dpb.FileDescriptorSet, error) {
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return fs, nil
}

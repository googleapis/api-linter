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
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/googleapis/api-linter/internal"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/proto"
	dpb "google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v3"
)

type cli struct {
	ConfigPath                string
	FormatType                string
	OutputPath                string
	ExitStatusOnLintFailure   bool
	VersionFlag               bool
	ProtoImportPaths          []string
	ProtoFiles                []string
	ProtoDescPath             []string
	EnabledRules              []string
	DisabledRules             []string
	ListRulesFlag             bool
	DebugFlag                 bool
	IgnoreCommentDisablesFlag bool
	RulePluginPaths           []string
}

// ExitForLintFailure indicates that a problem was found during linting.
//
//lint:ignore ST1012 modifying this variable name is a breaking change.
var ExitForLintFailure = errors.New("found problems during linting")

func newCli(args []string) *cli {
	// Define flag variables.
	var cfgFlag string
	var fmtFlag string
	var outFlag string
	var setExitStatusOnLintFailure bool
	var versionFlag bool
	var protoImportFlag []string
	var protoDescFlag []string
	var ruleEnableFlag []string
	var ruleDisableFlag []string
	var listRulesFlag bool
	var debugFlag bool
	var ignoreCommentDisablesFlag bool
	var rulePluginFlag []string

	// Register flag variables.
	fs := pflag.NewFlagSet("api-linter", pflag.ExitOnError)
	fs.StringVar(&cfgFlag, "config", "", "The linter config file.")
	fs.StringVar(&fmtFlag, "output-format", "", "The format of the linting results.\nSupported formats include \"yaml\", \"json\",\"github\" and \"summary\" table.\nYAML is the default.")
	fs.StringVarP(&outFlag, "output-path", "o", "", "The output file path.\nIf not given, the linting results will be printed out to STDOUT.")
	fs.BoolVar(&setExitStatusOnLintFailure, "set-exit-status", false, "Return exit status 1 when lint errors are found.")
	fs.BoolVar(&versionFlag, "version", false, "Print version and exit.")
	fs.StringArrayVarP(&protoImportFlag, "proto-path", "I", nil, "The folder for searching proto imports.\nMay be specified multiple times; directories will be searched in order.\nThe current working directory is always used.")
	fs.StringArrayVar(&protoDescFlag, "descriptor-set-in", nil, "The file containing a FileDescriptorSet for searching proto imports.\nMay be specified multiple times.")
	fs.StringArrayVar(&ruleEnableFlag, "enable-rule", nil, "Enable a rule with the given name.\nMay be specified multiple times.")
	fs.StringArrayVar(&ruleDisableFlag, "disable-rule", nil, "Disable a rule with the given name.\nMay be specified multiple times.")
	fs.BoolVar(&listRulesFlag, "list-rules", false, "Print the rules and exit.  Honors the output-format flag.")
	fs.BoolVar(&debugFlag, "debug", false, "Run in debug mode. Panics will print stack.")
	fs.BoolVar(&ignoreCommentDisablesFlag, "ignore-comment-disables", false, "If set to true, disable comments will be ignored.\nThis is helpful when strict enforcement of AIPs are necessary and\nproto definitions should not be able to disable checks.")
	fs.StringArrayVar(&rulePluginFlag, "rule-plugin", nil, "The path to a custom rule plugin (.so file).\nMay be specified multiple times.")

	// Parse flags.
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	return &cli{
		ConfigPath:                cfgFlag,
		FormatType:                fmtFlag,
		OutputPath:                outFlag,
		ExitStatusOnLintFailure:   setExitStatusOnLintFailure,
		ProtoImportPaths:          protoImportFlag,
		ProtoDescPath:             protoDescFlag,
		EnabledRules:              ruleEnableFlag,
		DisabledRules:             ruleDisableFlag,
		ProtoFiles:                fs.Args(),
		VersionFlag:               versionFlag,
		ListRulesFlag:             listRulesFlag,
		DebugFlag:                 debugFlag,
		IgnoreCommentDisablesFlag: ignoreCommentDisablesFlag,
		RulePluginPaths:           rulePluginFlag,
	}
}

func (c *cli) lint(rules lint.RuleRegistry, configs lint.Configs) error {
	// Print version and exit if asked.
	if c.VersionFlag {
		fmt.Printf("api-linter %s\n", internal.Version)
		return nil
	}

	if c.ListRulesFlag {
		return outputRules(c.FormatType)
	}

	// Pre-check if there are files to lint.
	if len(c.ProtoFiles) == 0 {
		return fmt.Errorf("no file to lint")
	}

	// Load custom rule plugins if provided
	if len(c.RulePluginPaths) > 0 {
		if err := loadCustomRulePlugins(c.RulePluginPaths, rules); err != nil {
			return fmt.Errorf("failed to load custom rule plugins: %v", err)
		}
		fmt.Printf("Loaded %d custom rule plugin(s)\n", len(c.RulePluginPaths))
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
	configs = append(configs, lint.Config{
		EnabledRules: c.EnabledRules,
	})
	// Add configs for the disabled rules.
	configs = append(configs, lint.Config{
		DisabledRules: c.DisabledRules,
	})
	// Prepare proto import lookup.
	fs, err := loadFileDescriptors(c.ProtoDescPath...)
	if err != nil {
		return err
	}
	lookupImport := func(name string) (*desc.FileDescriptor, error) {
		if f, found := fs[name]; found {
			return f, nil
		}
		return nil, fmt.Errorf("%q is not found", name)
	}
	var errorsWithPos []protoparse.ErrorWithPos
	var lock sync.Mutex
	// Parse proto files into `protoreflect` file descriptors.
	p := protoparse.Parser{
		ImportPaths:           append(c.ProtoImportPaths, "."),
		IncludeSourceCodeInfo: true,
		LookupImport:          lookupImport,
		ErrorReporter: func(errorWithPos protoparse.ErrorWithPos) error {
			// Protoparse isn't concurrent right now but just to be safe for the future.
			lock.Lock()
			errorsWithPos = append(errorsWithPos, errorWithPos)
			lock.Unlock()
			// Continue parsing. The error returned will be protoparse.ErrInvalidSource.
			return nil
		},
	}
	// Resolve file absolute paths to relative ones.
	// Using supplied import paths first.
	protoFiles := c.ProtoFiles
	if len(c.ProtoImportPaths) > 0 {
		protoFiles, err = protoparse.ResolveFilenames(c.ProtoImportPaths, c.ProtoFiles...)
		if err != nil {
			return err
		}
	}
	// Then resolve again against ".", the local directory.
	// This is necessary because ResolveFilenames won't resolve a path if it
	// relative to *at least one* of the given import paths, which can result
	// in duplicate file parsing and compilation errors, as seen in #1465 and
	// #1471. So we resolve against local (default) and flag specified import
	// paths separately.
	protoFiles, err = protoparse.ResolveFilenames([]string{"."}, protoFiles...)
	if err != nil {
		return err
	}
	fd, err := p.ParseFiles(protoFiles...)
	if err != nil {
		if err == protoparse.ErrInvalidSource {
			if len(errorsWithPos) == 0 {
				return errors.New("got protoparse.ErrInvalidSource but no ErrorWithPos errors")
			}
			// TODO: There's multiple ways to deal with this but this prints all the errors at least
			errStrings := make([]string, len(errorsWithPos))
			for i, errorWithPos := range errorsWithPos {
				errStrings[i] = errorWithPos.Error()
			}
			return errors.New(strings.Join(errStrings, "\n"))
		}
		return err
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rules, configs, lint.Debug(c.DebugFlag), lint.IgnoreCommentDisables(c.IgnoreCommentDisablesFlag))
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
	marshal := getOutputFormatFunc(c.FormatType)

	// Print the results.
	b, err := marshal(results)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}

	// Return error on lint failure which subsequently
	// exits with a non-zero status code
	if c.ExitStatusOnLintFailure && anyProblems(results) {
		return ExitForLintFailure
	}

	return nil
}

func anyProblems(results []lint.Response) bool {
	for i := range results {
		if len(results[i].Problems) > 0 {
			return true
		}
	}
	return false
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
	in, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return fs, nil
}

var outputFormatFuncs = map[string]formatFunc{
	"yaml": yaml.Marshal,
	"yml":  yaml.Marshal,
	"json": json.Marshal,
	"github": func(i interface{}) ([]byte, error) {
		switch v := i.(type) {
		case []lint.Response:
			return formatGitHubActionOutput(v), nil
		default:
			return json.Marshal(v)
		}
	},
	"summary": func(i interface{}) ([]byte, error) {
		switch v := i.(type) {
		case []lint.Response:
			return printSummaryTable(v)
		case listedRules:
			return v.printSummaryTable()
		default:
			return json.Marshal(v)
		}
	},
}

type formatFunc func(interface{}) ([]byte, error)

func getOutputFormatFunc(formatType string) formatFunc {
	if f, found := outputFormatFuncs[strings.ToLower(formatType)]; found {
		return f
	}
	return yaml.Marshal
}

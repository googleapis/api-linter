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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protocompile/reporter"
	"github.com/googleapis/api-linter/v2/internal"
	"github.com/googleapis/api-linter/v2/lint"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
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
	fs.BoolVar(&listRulesFlag, "list-rules", false, "Print the rules and exit. Honors the output-format flag.")
	fs.BoolVar(&debugFlag, "debug", false, "Run in debug mode. Panics will print stack.")
	fs.BoolVar(&ignoreCommentDisablesFlag, "ignore-comment-disables", false, "If set to true, disable comments will be ignored.\nThis is helpful when strict enforcement of AIPs are necessary and\nproto definitions should not be able to disable checks.")

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
	// Read linter config and append it to the default.
	if c.ConfigPath != "" {
		config, err := lint.ReadConfigsFromFile(c.ConfigPath)
		if err != nil {
			return err
		}
		configs = append(configs, config...)
	}
	// Add configs for the enabled and disabled rules from flags.
	// Combine them into a single config so that enable/disable
	// precedence is handled correctly.
	if len(c.EnabledRules) > 0 || len(c.DisabledRules) > 0 {
		configs = append(configs, lint.Config{
			EnabledRules:  c.EnabledRules,
			DisabledRules: c.DisabledRules,
		})
	}

	// Create resolver for descriptor sets.
	descResolver, err := loadFileDescriptorsAsResolver(c.ProtoDescPath...)
	if err != nil {
		return err
	}

	// Create resolver for source files.
	imports := resolveImports(c.ProtoImportPaths)
	sourceResolver := &protocompile.SourceResolver{
		ImportPaths: imports,
	}

	// This combines resolvers, prioritizing the source resolver and falling
	// back to the descriptor set resolver. This approach provides more accurate
	// descriptor information when the descriptor set lacks source details
	resolvers := []protocompile.Resolver{sourceResolver}
	if descResolver != nil {
		resolvers = append(resolvers, descResolver)
	}

	// The previous parser (`jhump/protoreflect`) reported all parse errors it
	// found. The default behavior of the new parser (`protocompile`) is to
	// stop on the first error.
	//
	// To preserve the original behavior, we provide a custom reporter that
	// collects all errors and allows the compiler to continue. The previous
	// parser also had no distinct concept of warnings, so we pass a nil
	// warning handler to maintain the same behavior of ignoring them.
	var collectedErrors []error
	rep := reporter.NewReporter(func(err reporter.ErrorWithPos) error {
		collectedErrors = append(collectedErrors, err)
		return nil // Returning nil signals the compiler to continue.
	}, nil)

	compiler := protocompile.Compiler{
		Resolver:       protocompile.WithStandardImports(protocompile.CompositeResolver(resolvers)),
		SourceInfoMode: protocompile.SourceInfoExtraOptionLocations,
		Reporter:       rep,
	}

	// Compile each file individually to avoid possible collisions
	// between a linted file that imports other files that are also being linted.
	// Otherwise, both the import resolver and the file will be "duplicated".
	var compiledFiles linker.Files
	for _, protoFile := range c.ProtoFiles {
		// The compiler returns a slice of files, even for a single input file.
		f, err := compiler.Compile(context.Background(), protoFile)
		// After compilation, check if the handler collected any errors.
		// This is the primary source of truth for parse errors when using a
		// custom reporter that continues on error.
		if len(collectedErrors) > 0 {
			errorStrings := make([]string, len(collectedErrors))
			for i, e := range collectedErrors {
				errorStrings[i] = e.Error()
			}
			return errors.New(strings.Join(errorStrings, "\n"))
		}

		// If the reporter has no errors, but the compiler still returned one,
		// it's a fatal, non-recoverable error.
		if err != nil {
			return err
		}
		// Append the compiled file(s) to the slice.
		compiledFiles = append(compiledFiles, f...)
	}
	files := compiledFiles

	// The compiler returns a slice of `*linker.File`, which is the compiler's
	// internal representation. We convert this to a slice of the standard
	// `protoreflect.FileDescriptor` interface, which the linter engine expects.
	var fileDescriptors []protoreflect.FileDescriptor
	for _, f := range files {
		fileDescriptors = append(fileDescriptors, f)
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rules, configs, lint.Debug(c.DebugFlag), lint.IgnoreCommentDisables(c.IgnoreCommentDisablesFlag))
	results, err := l.LintProtos(fileDescriptors...)
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

// resolver is a minimal implementation of the protocompile.Resolver interface.
// It is used to wrap a protoregistry.Files object, which is created from
// pre-compiled FileDescriptorSet files (`.protoset`), allowing the compiler
// to find and use these files for import resolution.
type resolver struct {
	files *protoregistry.Files
}

// FindFileByPath satisfies the protocompile.Resolver interface by searching
// for a file descriptor in the wrapped protoregistry.Files.
func (r *resolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	fd, err := r.files.FindFileByPath(path)
	if err != nil {
		return protocompile.SearchResult{}, err
	}
	return protocompile.SearchResult{Desc: fd}, nil
}

// loadFileDescriptorsAsResolver reads one or more FileDescriptorSet files
// (typically `.protoset` files) and loads them into a protoregistry.Files
// object. It then wraps this object in our custom resolver so that it can be
// used by the protocompile.Compiler to resolve imports.
func loadFileDescriptorsAsResolver(filePaths ...string) (protocompile.Resolver, error) {
	if len(filePaths) == 0 {
		return nil, nil
	}

	fdsSet := make(map[string]*dpb.FileDescriptorProto)
	for _, filePath := range filePaths {
		fs, err := readFileDescriptorSet(filePath)
		if err != nil {
			return nil, err
		}
		for _, fd := range fs.GetFile() {
			if _, exists := fdsSet[fd.GetName()]; !exists {
				fdsSet[fd.GetName()] = fd
			}
		}
	}

	fds := &dpb.FileDescriptorSet{}
	for _, fd := range fdsSet {
		fds.File = append(fds.File, fd)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, fmt.Errorf("failed to create protoregistry.Files: %w", err)
	}
	return &resolver{files: files}, nil
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

func resolveImports(imports []string) []string {
	// If no import paths are provided, default to the current directory.
	if len(imports) == 0 {
		return []string{"."}
	}

	// Get the absolute path of the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback: If we can't get CWD, return only the provided paths and "."
		seen := map[string]bool{
			".": true,
		}
		result := []string{"."} // Always include "."
		for _, p := range imports {
			if !seen[p] {
				seen[p] = true
				result = append(result, p)
			}
		}
		return result
	}

	// Resolve the canonical path for the current working directory.
	// This helps with symlinks (e.g., /var vs /private/var on macOS).
	evaluatedCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		// Fallback to Clean if EvalSymlinks fails (e.g., path does not exist)
		evaluatedCwd = filepath.Clean(cwd)
	}

	// Initialize resolvedImports with "." and track its canonical absolute path.
	resolvedImports := []string{"."}
	seenAbsolutePaths := map[string]bool{
		evaluatedCwd: true, // Mark canonical CWD as seen
	}

	for _, p := range imports {
		absPath, err := filepath.Abs(p)
		if err != nil {
			// If we can't get the absolute path, treat it as an external path
			// and add it if not already seen (by its original string form).
			if !seenAbsolutePaths[p] {
				seenAbsolutePaths[p] = true
				resolvedImports = append(resolvedImports, p)
			}
			continue
		}

		// Resolve the canonical path for the current import path.
		evaluatedAbsPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			// Fallback to Clean if EvalSymlinks fails
			evaluatedAbsPath = filepath.Clean(absPath)
		}

		// Check if the current import path's canonical form is the CWD's canonical form.
		// If so, it's covered by ".", so we skip it.
		if evaluatedAbsPath == evaluatedCwd {
			continue
		}

		// Add the original path if its canonical absolute form has not been seen before.
		if !seenAbsolutePaths[evaluatedAbsPath] {
			seenAbsolutePaths[evaluatedAbsPath] = true
			resolvedImports = append(resolvedImports, p)
		}
	}

	return resolvedImports
}

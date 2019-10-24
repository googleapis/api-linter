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

// The command line `api-lint` checks Google APIs defined in Protobuf files.
// It follows the API Improvement Proposals defined in https://aip.dev.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
	"github.com/jhump/protoreflect/desc/protoparse"
	"gopkg.in/yaml.v2"
)

var rules, _ = lint.NewRuleRegistry()
var configs lint.Configs

func init() {
	configs = lint.Configs{
		lint.Config{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"core": {},
			},
		},
	}

	if err := addRules(core.Rules()); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	if err := runCLI(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}

func runCLI(args []string) error {
	// Define flags.
	var fmtFlag string
	var cfgFlag string
	var outFlag string
	var protoImportFlag stringSlice

	// Register flags.
	fs := flag.NewFlagSet("api-lint", flag.ExitOnError)
	fs.Var(&protoImportFlag, "I", "The folders to search for proto imports.")
	fs.StringVar(&cfgFlag, "c", "", "The config file.")
	fs.StringVar(&fmtFlag, "f", "", "The format of the printed results.")
	fs.StringVar(&outFlag, "o", "", "The output file path.")

	// Parse flags.
	fs.Parse(args)
	// Add the current directory for searching imports.
	protoImportFlag = append(protoImportFlag, ".")
	files := fs.Args()
	// Check if there are files to lint; if not, abort.
	if len(files) == 0 {
		return fmt.Errorf("no files to lint")
	}

	// Parse files into protoreflect file descriptors.
	p := protoparse.Parser{
		ImportPaths:           protoImportFlag,
		IncludeSourceCodeInfo: true,
	}
	fd, err := p.ParseFiles(files...)
	if err != nil {
		return err
	}

	// Parse the provided config.
	if cfgFlag != "" {
		c, err := lint.ReadConfigsFromFile(cfgFlag)
		if err != nil {
			return err
		}
		configs = append(configs, c...)
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rules, configs)
	results, err := l.LintProtos(fd...)
	if err != nil {
		return err
	}

	// Determine the output to write the results.
	w := os.Stdout
	if outFlag != "" {
		var err error
		w, err = os.Create(outFlag)
		if err != nil {
			return err
		}
		defer w.Close()
	}

	// Determine the format to print the results.
	marshal := yaml.Marshal
	switch fmtFlag {
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

func addRules(r lint.RuleRegistry) error {
	return rules.Register(r.All()...)
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

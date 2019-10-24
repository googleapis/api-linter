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

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func runCLI(rules lint.RuleRegistry, configs lint.Configs, args []string) error {
	app := cli.NewApp()
	app.Name = "api-linter"
	app.Usage = "A linter for APIs"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:      "check",
			Aliases:   []string{"c"},
			Usage:     "Check protobuf files that define an API",
			ArgsUsage: "files...",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cfg",
					Value: "",
					Usage: "configuration file path",
				},
				cli.StringFlag{
					Name:  "out",
					Value: "",
					Usage: "output file path (default: stdout)",
				},
				cli.StringFlag{
					Name:  "fmt",
					Value: "yaml",
					Usage: "output format",
				},
				cli.StringSliceFlag{
					Name:  "proto_path",
					Value: &cli.StringSlice{"."},
					Usage: "the directories in which for proto parser to search for imports",
				},
				cli.StringFlag{
					Name:  "proto_desc",
					Value: "",
					Usage: "the compiled file descriptor set for searching imports",
				},
			},
			Action: func(c *cli.Context) error {
				filenames := c.Args()

				// Sanity check: Were we given any files to parse at all?
				// If not, abort.
				if len(filenames) == 0 {
					os.Stderr.WriteString("No files specified to lint.\n")
					return nil
				}

				// Prepare lookup for proto imports, in additional to the proto paths.
				var lookupImport func(string) (*desc.FileDescriptor, error)
				if c.String("proto_desc") != "" {
					descriptors, err := loadFileDescriptors(c.String("proto_desc"))
					if err != nil {
						return err
					}
					lookupImport = func(name string) (*desc.FileDescriptor, error) {
						if d, found := descriptors[name]; found {
							return d, nil
						}
						return nil, fmt.Errorf("%q is not found in the file %q", name, c.String("proto_desc"))
					}
				}

				// Parse the provided protobuf files into a protoreflect file
				// descriptor.
				p := protoparse.Parser{
					ImportPaths:           c.StringSlice("proto_path"),
					IncludeSourceCodeInfo: true,
					LookupImport:          lookupImport,
				}
				fd, err := p.ParseFiles(filenames...)
				if err != nil {
					return err
				}

				// If a configuration file was provided, parse it.
				if c.String("cfg") != "" {
					userConfigs, err := lint.ReadConfigsFromFile(c.String("cfg"))
					if err != nil {
						return err
					}
					configs = append(configs, userConfigs...)
				}

				// Instantiate the linter object, and lint the protos.
				l := lint.New(rules, configs)
				lintResponses, err := l.LintProtos(fd...)
				if err != nil {
					return err
				}

				// If writing output to a file, set that up.
				// If no file was specified, use stdout.
				w := os.Stdout
				if c.String("out") != "" {
					var err error
					w, err = os.Create(c.String("out"))
					if err != nil {
						return err
					}
					defer w.Close()
				}

				// Determine what format we are using to print the results.
				marshal := yaml.Marshal
				switch c.String("fmt") {
				case "json":
					marshal = json.Marshal
				case "summary":
					marshal = func(i interface{}) ([]byte, error) {
						return emitSummary(i.([]lint.Response))
					}
				}

				// Print the actual results.
				b, err := marshal(lintResponses)
				if err != nil {
					return err
				}
				if _, err = w.Write(b); err != nil {
					return err
				}
				return nil
			},
		},
	}

	return app.Run(args)
}

func loadFileDescriptors(filename string) (map[string]*desc.FileDescriptor, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorsFromSet(fs)
}

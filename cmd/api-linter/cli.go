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
	"os"

	"github.com/googleapis/api-linter/lint"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func runCLI(rules lint.Rules, configs lint.Configs, args []string) error {
	app := cli.NewApp()
	app.Name = "api-linter"
	app.Usage = "A linter for APIs"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:      "checkproto",
			Aliases:   []string{"cp"},
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
				cli.StringFlag{
					Name:  "protoc",
					Value: "protoc",
					Usage: "protocol compiler path",
				},
				cli.StringFlag{
					Name:  "proto_path",
					Value: ".",
					Usage: "the directory in which for protoc to search for imports",
				},
			},
			Action: func(c *cli.Context) error {
				filenames := c.Args()
				if len(filenames) == 0 {
					return nil
				}

				p := protocParser{
					importPath: c.String("proto_path"),
					protoc:     c.String("protoc"),
				}
				files, err := p.ParseProto(filenames...)
				if err != nil {
					return err
				}

				if c.String("cfg") != "" {
					userConfigs, err := lint.ReadConfigsFromFile(c.String("cfg"))
					if err != nil {
						return err
					}
					configs = append(configs, userConfigs...)
				}

				l := lint.New(rules, configs)
				lintResponses, err := l.LintProtos(files)
				if err != nil {
					return err
				}

				w := os.Stdout
				if c.String("out") != "" {
					var err error
					w, err = os.Create(c.String("out"))
					if err != nil {
						return err
					}
					defer w.Close()
				}

				marshal := yaml.Marshal
				switch c.String("fmt") {
				case "json":
					marshal = json.Marshal
				}

				if c.String("fmt") == "summary" {
					summary := createSummary(lintResponses)
					emitSummary(&summary, w)
				} else {
					b, err := marshal(lintResponses)
					if err != nil {
						return err
					}
					if _, err = w.Write(b); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	return app.Run(args)
}

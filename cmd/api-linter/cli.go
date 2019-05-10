package main

import (
	"encoding/json"
	"os"

	"github.com/googleapis/api-linter/lint"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func runCLI(rules []lint.Rule, configs lint.RuntimeConfigs, args []string) error {
	app := cli.NewApp()
	app.Name = "api-linter"
	app.Usage = "A linter for APIs defined in protocol buffers."
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

				l := newLinter(rules, configs)
				problems, err := l.LintProto(files)
				if err != nil {
					return err
				}

				outFile := c.String("out")
				outFmt := c.String("fmt")
				return print(problems, outFile, outFmt)
			},
		},
		{
			Name:      "print",
			Aliases:   []string{"p"},
			Usage:     "Print rule information",
			ArgsUsage: "rules...",
			Flags: []cli.Flag{
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
			},
			Action: func(c *cli.Context) error {
				var infors []lint.RuleInfo
				for _, rl := range rules {
					infors = append(infors, rl.Info())
				}

				if len(c.Args()) > 0 {
					prefixes := c.Args()
					var filtered []lint.RuleInfo
					for _, ri := range infors {
						for _, pre := range prefixes {
							if ri.Name.HasPrefix(pre) {
								filtered = append(filtered, ri)
								break
							}
						}
					}
					infors = filtered
				}

				outFile := c.String("out")
				outFmt := c.String("fmt")
				return print(infors, outFile, outFmt)
			},
		},
	}

	return app.Run(args)
}

// print data to the file in the given format. By default, it prints
// to stdout if file is empty, and in YAML format if fmt is empty.
func print(data interface{}, file, fmt string) error {
	w := os.Stdout
	if file != "" {
		var err error
		w, err = os.Create(file)
		if err != nil {
			return err
		}
		defer w.Close()
	}
	marshal := yaml.Marshal
	switch fmt {
	case "json":
		marshal = json.Marshal
	}
	b, err := marshal(data)
	if err != nil {
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

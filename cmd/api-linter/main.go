package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/googleapis/api-linter/lint"
	core "github.com/googleapis/api-linter/rules"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fmt",
			Value: "yaml",
			Usage: "output format",
		},
		cli.StringFlag{
			Name:  "cfg",
			Value: "",
			Usage: "configuration file",
		},
		cli.StringFlag{
			Name:  "import_path",
			Value: ".",
			Usage: "protoc import path",
		},
		cli.StringFlag{
			Name:  "out",
			Value: "",
			Usage: "output file",
		},
	}
	app.Action = func(c *cli.Context) error {
		filenames := c.Args()
		if len(filenames) == 0 {
			return nil
		}

		p := protocParser{
			importPath: c.String("import_path"),
		}
		files, err := p.ParseProto(filenames...)
		if err != nil {
			return err
		}

		l := newLinter(getRules(), getConfigs())
		problems, err := l.LintProto(files)
		if err != nil {
			return err
		}

		w := os.Stderr
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
		b, err := marshal(problems)
		if err != nil {
			return err
		}

		if _, err = w.Write(b); err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Register default configuration.
func getConfigs() lint.RuntimeConfigs {
	return lint.RuntimeConfigs{
		lint.RuntimeConfig{
			IncludedPaths: []string{"**/*.proto"},
			RuleConfigs: map[string]lint.RuleConfig{
				"core": {
					Status:   lint.Enabled,
					Category: lint.Warning,
				},
			},
		},
	}
}

// Register rules.
func getRules() []lint.Rule {
	var rules []lint.Rule
	rules = append(rules, core.Rules().All()...)
	return rules
}

// Encoder is an interface wrapping around Encode method.
type Encoder interface {
	Encode(interface{}) error
}

func writeProblems(problems []lint.Problem, e Encoder) error {
	return e.Encode(problems)
}

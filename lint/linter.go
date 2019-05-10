package lint

import (
	"errors"
	"strings"

	descriptorpb "github.com/golang/protobuf/v2/types/descriptor"
)

// Linter checks API files and returns a list of detected problems.
type Linter struct {
	rules   Rules
	configs Configs
}

// New creates and returns a linter with the given rules and configs.
func New(rules Rules, configs Configs) *Linter {
	l := &Linter{
		rules:   rules,
		configs: configs,
	}
	return l
}

// LintProtos checks protobuf files and returns a list of problems or an error.
func (l *Linter) LintProtos(files []*descriptorpb.FileDescriptorProto) ([]Response, error) {
	return l.lintProtos(files)
}

func (l *Linter) lintProtos(files []*descriptorpb.FileDescriptorProto) ([]Response, error) {
	var responses []Response
	for _, proto := range files {
		req, err := NewProtoRequest(proto)
		if err != nil {
			return nil, err
		}
		resp, err := l.run(req)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

// run executes rules on the request.
//
// It uses the proto file path to determine which rules will
// be applied to the request, according to the list of Linter
// configs.
func (l *Linter) run(req Request) (Response, error) {
	resp := Response{
		FilePath: req.ProtoFile().Path(),
	}
	var errMessages []string

	for name, rl := range l.rules {
		config := getDefaultRuleConfig()

		if c, err := l.configs.getRuleConfig(req.ProtoFile().Path(), name); err == nil {
			config = config.withOverride(c)
		} else {
			errMessages = append(errMessages, err.Error())
			continue
		}

		if config.Status == Enabled && !req.DescriptorSource().isRuleDisabledInFile(rl.Info().Name) {
			if problems, err := rl.Lint(req); err == nil {
				for _, p := range problems {
					if !req.DescriptorSource().isRuleDisabled(rl.Info().Name, p.Descriptor) {
						p.RuleID = rl.Info().Name
						p.Category = config.Category
						resp.Problems = append(resp.Problems, p)
					}
				}
			} else {
				errMessages = append(errMessages, err.Error())
			}
		}
	}

	var err error
	if len(errMessages) != 0 {
		err = errors.New(strings.Join(errMessages, "; "))
	}

	return resp, err
}

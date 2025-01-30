package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checkutil"
	"buf.build/go/bufplugin/descriptor"
	"buf.build/go/bufplugin/info"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var allCategory = "GOOGLE_ALL"

func main() {
	spec, err := buildSpec()
	if err != nil {
		log.Fatalf("failed to build buf plugin spec: %v", err)
	}
	check.Main(spec)
}

// buildSpec creates the buf plugin config
func buildSpec() (*check.Spec, error) {
	spec := &check.Spec{
		Rules: []*check.RuleSpec{},
		Categories: []*check.CategorySpec{
			{ID: allCategory, Purpose: "Checks all Google api-linter rules."},
		},
		Info: &info.Spec{
			Documentation: `A linting plugin that checks Google AIP conformance via the api-linter project`,
			SPDXLicenseID: "apache-2.0",
			LicenseURL:    "https://github.com/googleapis/api-linter/blob/main/LICENSE",
		},
	}
	// add all rules. buf.yaml specifies which to enable/disable via ID
	ruleRegistry := lint.NewRuleRegistry()
	if err := rules.Add(ruleRegistry); err != nil {
		return nil, fmt.Errorf("error when registering lint rules: %w", err)
	}
	for _, rule := range ruleRegistry {
		spec.Rules = append(spec.Rules, ruleToBufRule(rule))
	}
	return spec, nil
}

// ruleNameToBufID sanitizes the rule names into a format accepted by buf
func ruleNameToBufID(name lint.RuleName) string {
	replacer := strings.NewReplacer("::", "_", "-", "_")
	return "GOOGLE_" + strings.ToUpper(replacer.Replace(string(name)))
}

// ruleToBufRule translate's our "rule" into the buf equivalant
func ruleToBufRule(rule lint.ProtoRule) *check.RuleSpec {
	handler := func(ctx context.Context, writer check.ResponseWriter, request check.Request, descriptor descriptor.FileDescriptor) error {
		wrappedDescriptor, err := desc.WrapFile(descriptor.ProtoreflectFileDescriptor())
		if err != nil {
			return err
		}
		problems := rule.Lint(wrappedDescriptor)
		for _, problem := range problems {
			var unwrapedDescriptor protoreflect.Descriptor
			switch d := problem.Descriptor.(type) {
			case *desc.FieldDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.FileDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.EnumValueDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.EnumDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.MessageDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.MethodDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.OneOfDescriptor:
				unwrapedDescriptor = d.Unwrap()
			case *desc.ServiceDescriptor:
				unwrapedDescriptor = d.Unwrap()
			default:
				return fmt.Errorf("unhandled type converting api-linter rule to buf rule: %T", d)
			}
			writer.AddAnnotation(
				check.WithMessage(problem.Message),
				check.WithDescriptor(unwrapedDescriptor),
			)
		}
		return nil
	}
	return &check.RuleSpec{
		ID:          ruleNameToBufID(rule.GetName()),
		CategoryIDs: []string{allCategory},
		Type:        check.RuleTypeLint,
		Purpose:     fmt.Sprintf("Validates Google AIP %s.", rule.GetName()),
		Default:     true,
		Handler:     checkutil.NewFileRuleHandler(handler, checkutil.WithoutImports()),
	}
}

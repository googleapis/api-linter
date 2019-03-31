// Package rules contains rules that checks API styles
// in Google Protobuf APIs.
package rules

import (
	"log"
	"strings"

	"github.com/golang/protobuf/reflect/protoreflect"
	"github.com/jgeewax/api-linter/lint"
	"github.com/jgeewax/api-linter/visitors"
)

var coreRules = lint.NewRules()

func registerRuleWithLintFunc(metadata ruleMetadata, lintFunc lintFuncType) {
	r := ruleBase{
		metadata: metadata,
		lintFunc: lintFunc,
	}
	registerRuleBase(r)
}

func registerRuleWithVisitor(metadata ruleMetadata, visitor protoRuleVisitor) {
	r := ruleBase{
		metadata: metadata,
		visitor:  visitor,
	}
	registerRuleBase(r)
}

func registerRuleBase(r ruleBase) {
	if err := coreRules.Register(r); err != nil {
		log.Fatalf("Error when registering rule '%s': %v", r.ID(), err)
	}
}

// Rules returns all rules registered in this package.
func Rules() *lint.Rules {
	return coreRules
}

// ruleMetadata stores metadata information of a lint.Rule.
type ruleMetadata struct {
	Set         string          // rule set name.
	Name        string          // rule name in the set.
	Description string          // a short description of this rule.
	URL         string          // a link to a document for more details.
	FileTypes   []lint.FileType // types of files that this rule targets to.
	Category    lint.Category   // category of problems this rule produces.
}

type lintFuncType func(lint.Request) (lint.Response, error)

// ruleBase implements lint.Rule.
type ruleBase struct {
	metadata ruleMetadata
	lintFunc lintFuncType
	visitor  protoRuleVisitor
}

func (r ruleBase) ID() lint.RuleID {
	return lint.RuleID{
		Set:  r.metadata.Set,
		Name: r.metadata.Name,
	}
}

func (r ruleBase) Description() string {
	return r.metadata.Description
}

func (r ruleBase) URL() string {
	return r.metadata.URL
}

func (r ruleBase) FileTypes() []lint.FileType {
	return r.metadata.FileTypes
}

func (r ruleBase) Category() lint.Category {
	return r.metadata.Category
}

func (r ruleBase) Lint(req lint.Request) (lint.Response, error) {
	if r.lintFunc != nil {
		return r.lintFunc(req)
	}
	if r.visitor != nil {
		return r.visitor.Walk(r, req)
	}
	return lint.Response{}, nil
}

type protoRuleVisitor interface {
	Walk(r lint.Rule, req lint.Request) (lint.Response, error)
}

// simpleVisitor can check top level enums, all messages and services.
type simpleVisitor struct {
	MessageCheck   func(protoreflect.MessageDescriptor, lint.Context) []lint.Problem
	FieldCheck     func(protoreflect.FieldDescriptor, lint.Context) []lint.Problem
	EnumCheck      func(protoreflect.EnumDescriptor, lint.Context) []lint.Problem
	EnumValueCheck func(protoreflect.EnumValueDescriptor, lint.Context) []lint.Problem
	OneofCheck     func(protoreflect.OneofDescriptor, lint.Context) []lint.Problem
	ServiceCheck   func(protoreflect.ServiceDescriptor, lint.Context) []lint.Problem
	MethodCheck    func(protoreflect.MethodDescriptor, lint.Context) []lint.Problem

	rule     lint.Rule
	ctx      lint.Context
	problems []lint.Problem
}

func (v *simpleVisitor) Walk(r lint.Rule, req lint.Request) (lint.Response, error) {
	v.rule = r
	v.ctx = req.Context()
	v.checkAllMessages(req.ProtoFile())
	v.checkTopLevelEnums(req.ProtoFile())
	v.checkAllServices(req.ProtoFile())
	return lint.Response{Problems: v.problems}, nil
}

func (v *simpleVisitor) checkMessage(m protoreflect.MessageDescriptor) {
	if v.MessageCheck != nil && v.isRuleEnabled(m) {
		v.addProblems(v.MessageCheck(m, v.ctx)...)
	}
}

func (v *simpleVisitor) checkField(f protoreflect.FieldDescriptor) {
	if v.FieldCheck != nil && v.isRuleEnabled(f) {
		v.addProblems(v.FieldCheck(f, v.ctx)...)
	}
}

func (v *simpleVisitor) checkEnum(e protoreflect.EnumDescriptor) {
	if v.EnumCheck != nil && v.isRuleEnabled(e) {
		v.addProblems(v.EnumCheck(e, v.ctx)...)
	}
}

func (v *simpleVisitor) checkEnumValue(ev protoreflect.EnumValueDescriptor) {
	if v.EnumValueCheck != nil && v.isRuleEnabled(ev) {
		v.addProblems(v.EnumValueCheck(ev, v.ctx)...)
	}
}

func (v *simpleVisitor) checkOneof(o protoreflect.OneofDescriptor) {
	if v.OneofCheck != nil && v.isRuleEnabled(o) {
		v.addProblems(v.OneofCheck(o, v.ctx)...)
	}
}

func (v *simpleVisitor) checkTopLevelEnums(f protoreflect.FileDescriptor) {
	for i := 0; i < f.Enums().Len(); i++ {
		v.checkEnum(f.Enums().Get(i))
	}
}

func (v *simpleVisitor) checkAllMessages(f protoreflect.FileDescriptor) {
	visitors.WalkMessage(f, &visitors.SimpleMessageVisitor{
		Funcs: visitors.MessageVisitingFuncs{
			MessageVisit:   v.checkMessage,
			FieldVisit:     v.checkField,
			EnumVisit:      v.checkEnum,
			EnumValueVisit: v.checkEnumValue,
			OneofVisit:     v.checkOneof,
		},
	})
}

func (v *simpleVisitor) checkAllServices(f protoreflect.FileDescriptor) {
	visitors.WalkService(f, &visitors.SimpleServiceVisitor{
		Funcs: visitors.ServiceVisitFuncs{
			MethodVisit:  v.checkMethod,
			ServiceVisit: v.checkService,
		},
	})
}

func (v *simpleVisitor) checkService(s protoreflect.ServiceDescriptor) {
	if v.ServiceCheck != nil && v.isRuleEnabled(s) {
		v.addProblems(v.ServiceCheck(s, v.ctx)...)
	}
}

func (v *simpleVisitor) checkMethod(m protoreflect.MethodDescriptor) {
	if v.MethodCheck != nil && v.isRuleEnabled(m) {
		v.addProblems(v.MethodCheck(m, v.ctx)...)
	}
}

func (v *simpleVisitor) addProblems(problems ...lint.Problem) {
	v.problems = append(v.problems, problems...)
}

func (v *simpleVisitor) isRuleEnabled(d protoreflect.Descriptor) bool {
	return isRuleEnabled(v.rule.ID(), d, v.ctx)
}

func isRuleEnabled(ruleID lint.RuleID, d protoreflect.Descriptor, ctx lint.Context) bool {
	comments, err := ctx.DescriptorSource().FindCommentsByDescriptor(d)
	if err != nil {
		log.Printf("FindCommentsByDescriptor for '%s' returned error: %v", d.FullName(), err)
		return true
	}

	leadingAndInLineComments := []string{comments.LeadingComments, comments.TrailingComments}
	return !stringsContains(leadingAndInLineComments, ruleDisablingComment(ruleID))
}

func stringsContains(comments []string, s string) bool {
	for _, c := range comments {
		if strings.Contains(c, s) {
			return true
		}
	}
	return false
}

func ruleDisablingComment(id lint.RuleID) string {
	name := id.Set + "." + id.Name
	if id.Set == "" || id.Set == "core" {
		name = id.Name
	}
	return "(-- api-linter: " + name + "=disabled --)"
}

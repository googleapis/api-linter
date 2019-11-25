package descrule

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

func TestFileRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"FileDescriptor_LintInvoked", NewFile(&desc.FileDescriptor{}), true},
		{"NotAFileDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &FileRule{
				LintFile: func(f *desc.FileDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestFileRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.FileDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.FileDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.FileDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &FileRule{
				OnlyIf: test.onlyIf,
				LintFile: func(f *desc.FileDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewFile(&desc.FileDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestFileRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &FileRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestMessageRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"MessageDescriptor_LintInvoked", NewMessage(&desc.MessageDescriptor{}), true},
		{"NotAMessageDescriptor_LintNotInvoked", NewFile(&desc.FileDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &MessageRule{
				LintMessage: func(f *desc.MessageDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestMessageRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.MessageDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.MessageDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.MessageDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &MessageRule{
				OnlyIf: test.onlyIf,
				LintMessage: func(f *desc.MessageDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewMessage(&desc.MessageDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestMessageRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &MessageRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestFieldRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"FieldDescriptor_LintInvoked", NewField(&desc.FieldDescriptor{}), true},
		{"NotAFieldDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &FieldRule{
				LintField: func(f *desc.FieldDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestFieldRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.FieldDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.FieldDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.FieldDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &FieldRule{
				OnlyIf: test.onlyIf,
				LintField: func(f *desc.FieldDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewField(&desc.FieldDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestFieldRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &FieldRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestEnumRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"EnumDescriptor_LintInvoked", NewEnum(&desc.EnumDescriptor{}), true},
		{"NotAEnumDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &EnumRule{
				LintEnum: func(f *desc.EnumDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestEnumRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.EnumDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.EnumDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.EnumDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &EnumRule{
				OnlyIf: test.onlyIf,
				LintEnum: func(f *desc.EnumDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewEnum(&desc.EnumDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestEnumRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &EnumRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestEnumValueRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"EnumValueDescriptor_LintInvoked", NewEnumValue(&desc.EnumValueDescriptor{}), true},
		{"NotAEnumValueDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &EnumValueRule{
				LintEnumValue: func(f *desc.EnumValueDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestEnumValueRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.EnumValueDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.EnumValueDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.EnumValueDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &EnumValueRule{
				OnlyIf: test.onlyIf,
				LintEnumValue: func(f *desc.EnumValueDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewEnumValue(&desc.EnumValueDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestEnumValueRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &EnumValueRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestMethodRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"MethodDescriptor_LintInvoked", NewMethod(&desc.MethodDescriptor{}), true},
		{"NotAMethodDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &MethodRule{
				LintMethod: func(f *desc.MethodDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestMethodRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.MethodDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.MethodDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.MethodDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &MethodRule{
				OnlyIf: test.onlyIf,
				LintMethod: func(f *desc.MethodDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewMethod(&desc.MethodDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestMethodRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &MethodRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestServiceRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"ServiceDescriptor_LintInvoked", NewService(&desc.ServiceDescriptor{}), true},
		{"NotAServiceDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &ServiceRule{
				LintService: func(f *desc.ServiceDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestServiceRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(*desc.ServiceDescriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.ServiceDescriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(*desc.ServiceDescriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &ServiceRule{
				OnlyIf: test.onlyIf,
				LintService: func(f *desc.ServiceDescriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewService(&desc.ServiceDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestServiceRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &ServiceRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

func TestDescriptorRule_Lint(t *testing.T) {
	tests := []struct {
		name        string
		descriptor  lint.Descriptor
		lintInvoked bool
	}{
		{"DescriptorDescriptor_LintInvoked", NewDescriptor(&desc.FileDescriptor{}), true},
		{"NotADescriptorDescriptor_LintNotInvoked", NewMessage(&desc.MessageDescriptor{}), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &DescriptorRule{
				LintDescriptor: func(f desc.Descriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(test.descriptor)
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestDescriptorRule_OnlyIf(t *testing.T) {
	tests := []struct {
		name        string
		onlyIf      func(desc.Descriptor) bool
		lintInvoked bool
	}{
		{"NoOnlyIf_Invoked", nil, true},
		{"OnlyIf_ReturnFalse_NotInvoked", func(desc.Descriptor) bool { return false }, false},
		{"OnlyIf_ReturnTrue_Invoked", func(desc.Descriptor) bool { return true }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &DescriptorRule{
				OnlyIf: test.onlyIf,
				LintDescriptor: func(f desc.Descriptor) []lint.Problem {
					return []lint.Problem{}
				},
			}
			problems := rule.Lint(NewDescriptor(&desc.FileDescriptor{}))
			lintInvoked := problems != nil
			if test.lintInvoked != lintInvoked {
				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
			}
		})
	}
}

func TestDescriptorRule_GetName(t *testing.T) {
	want := lint.RuleName("test")
	rule := &DescriptorRule{
		RuleName: want,
	}
	if rule.Name() != want {
		t.Errorf("GetName() got %s, but want %s", rule.Name(), want)
	}
}

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

package lint

// func TestFileRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"FileDescriptor_LintInvoked", &desc.FileDescriptor{}, true},
// 		{"NotAFileDescriptor_LintNotInvoked", &desc.MessageDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &FileRule{
// 				LintFile: func(f *desc.FileDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestFileRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.FileDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.FileDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.FileDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &FileRule{
// 				OnlyIf: test.onlyIf,
// 				LintFile: func(f *desc.FileDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.FileDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestFileRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &FileRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

// func TestMessageRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"MessageDescriptor_LintInvoked", &desc.MessageDescriptor{}, true},
// 		{"NotAMessageDescriptor_LintNotInvoked", &desc.FileDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &MessageRule{
// 				LintMessage: func(f *desc.MessageDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestMessageRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.MessageDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.MessageDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.MessageDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &MessageRule{
// 				OnlyIf: test.onlyIf,
// 				LintMessage: func(f *desc.MessageDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.MessageDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestMessageRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &MessageRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

// func TestFieldRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"FieldDescriptor_LintInvoked", &desc.FieldDescriptor{}, true},
// 		{"NotAFieldDescriptor_LintNotInvoked", &desc.MessageDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &FieldRule{
// 				LintField: func(f *desc.FieldDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestFieldRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.FieldDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.FieldDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.FieldDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &FieldRule{
// 				OnlyIf: test.onlyIf,
// 				LintField: func(f *desc.FieldDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.FieldDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestFieldRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &FieldRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

// func TestEnumRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"EnumDescriptor_LintInvoked", &desc.EnumDescriptor{}, true},
// 		{"NotAEnumDescriptor_LintNotInvoked", &desc.MessageDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &EnumRule{
// 				LintEnum: func(f *desc.EnumDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestEnumRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.EnumDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.EnumDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.EnumDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &EnumRule{
// 				OnlyIf: test.onlyIf,
// 				LintEnum: func(f *desc.EnumDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.EnumDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestEnumRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &EnumRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

// func TestMethodRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"MethodDescriptor_LintInvoked", &desc.MethodDescriptor{}, true},
// 		{"NotAMethodDescriptor_LintNotInvoked", &desc.MessageDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &MethodRule{
// 				LintMethod: func(f *desc.MethodDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestMethodRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.MethodDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.MethodDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.MethodDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &MethodRule{
// 				OnlyIf: test.onlyIf,
// 				LintMethod: func(f *desc.MethodDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.MethodDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestMethodRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &MethodRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

// func TestServiceRule_Lint(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		descriptor  desc.Descriptor
// 		lintInvoked bool
// 	}{
// 		{"ServiceDescriptor_LintInvoked", &desc.ServiceDescriptor{}, true},
// 		{"NotAServiceDescriptor_LintNotInvoked", &desc.MessageDescriptor{}, false},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &ServiceRule{
// 				LintService: func(f *desc.ServiceDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(test.descriptor)
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestServiceRule_OnlyIf(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		onlyIf      func(*desc.ServiceDescriptor) bool
// 		lintInvoked bool
// 	}{
// 		{"NoOnlyIf_Invoked", nil, true},
// 		{"OnlyIf_ReturnFalse_NotInvoked", func(*desc.ServiceDescriptor) bool { return false }, false},
// 		{"OnlyIf_ReturnTrue_Invoked", func(*desc.ServiceDescriptor) bool { return true }, true},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			rule := &ServiceRule{
// 				OnlyIf: test.onlyIf,
// 				LintService: func(f *desc.ServiceDescriptor) []Problem {
// 					return []Problem{}
// 				},
// 			}
// 			problems := rule.Lint(&desc.ServiceDescriptor{})
// 			lintInvoked := problems != nil
// 			if test.lintInvoked != lintInvoked {
// 				t.Errorf("Lint() invoked? got %v, but want %v", lintInvoked, test.lintInvoked)
// 			}
// 		})
// 	}
// }

// func TestServiceRule_GetName(t *testing.T) {
// 	want := RuleName("test")
// 	rule := &ServiceRule{
// 		Name: want,
// 	}
// 	if rule.GetName() != want {
// 		t.Errorf("GetName() got %s, but want %s", rule.GetName(), want)
// 	}
// }

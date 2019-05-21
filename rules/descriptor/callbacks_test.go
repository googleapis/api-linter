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

package descriptor

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"github.com/googleapis/api-linter/lint"
)

func TestCallbacks_Apply(t *testing.T) {
	f := readProtoFile("test.protoset")
	descriptors := map[string]protoreflect.Descriptor{
		"enum":      f.Enums().Get(0),
		"enumvalue": f.Enums().Get(0).Values().Get(0),
		"field":     f.Messages().Get(0).Fields().Get(0),
		"message":   f.Messages().Get(0),
		"method":    f.Services().Get(0).Methods().Get(0),
		"oneof":     f.Messages().Get(0).Oneofs().Get(0),
		"service":   f.Services().Get(0),
	}

	tests := []struct {
		callbacks Callbacks
		results   []lint.Problem
	}{
		{Callbacks{}, []lint.Problem{}},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "enum"}}, nil
				},
			},
			[]lint.Problem{{Message: "enum"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "enumvalue"}}, nil
				},
			},
			[]lint.Problem{{Message: "enumvalue"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "field"}}, nil
				},
			},
			[]lint.Problem{{Message: "field"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "message"}}, nil
				},
			},
			[]lint.Problem{{Message: "message"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "method"}}, nil
				},
			},
			[]lint.Problem{{Message: "method"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "oneof"}}, nil
				},
			},
			[]lint.Problem{{Message: "oneof"}},
		},
		{
			Callbacks{
				EnumCallback: func(d protoreflect.EnumDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "service"}}, nil
				},
			},
			[]lint.Problem{{Message: "service"}},
		},
	}

	for _, test := range tests {
		results := []lint.Problem{}
		for _, d := range descriptors {
			problems, err := test.callbacks.Apply(d, lint.DescriptorSource{})
			if err != nil {
				t.Errorf("Callbacks.Apply returns unexpected error: %v", err)
			}
			results = append(results, problems...)
		}
		if got, want := results, test.results; !reflect.DeepEqual(got, want) {
			t.Errorf("Callbacks.Apply returns problems '%s', but want '%s'", got, want)
		}
	}
}

func TestCallbacks_Apply_DescriptorCallback(t *testing.T) {
	f := readProtoFile("test.protoset")
	descriptors := map[string]protoreflect.Descriptor{
		"enum":      f.Enums().Get(0),
		"enumvalue": f.Enums().Get(0).Values().Get(0),
	}

	tests := []struct {
		callbacks  Callbacks
		numProblem int
	}{
		{Callbacks{}, 0},
		{
			Callbacks{
				DescriptorCallback: func(d protoreflect.Descriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
					return []lint.Problem{{Message: "desc"}}, nil
				},
			},
			2,
		},
	}

	for _, test := range tests {
		all := []lint.Problem{}
		for _, d := range descriptors {
			problems, err := test.callbacks.Apply(d, lint.DescriptorSource{})
			if err != nil {
				t.Errorf("Callbacks.Apply returns unexpected error: %v", err)
			}
			all = append(all, problems...)
		}
		if got, want := len(all), test.numProblem; got != want {
			t.Errorf("Callbacks.Apply returns %d problems, but want %d", got, want)
		}
	}
}

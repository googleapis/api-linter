// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestGetHTTPRules(t *testing.T) {
	// Create a MethodOptions with the annotation set.
	opts := &dpb.MethodOptions{}
	httpRule := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{name=publishers/*/books/*}",
		},
		AdditionalBindings: []*annotations.HttpRule{{
			Pattern: &annotations.HttpRule_Get{
				Get: "/v1/{name=books/*}",
			},
		}},
	}

	if err := proto.SetExtension(opts, annotations.E_Http, httpRule); err != nil {
		t.Fatalf("Failed to set google.api.http annotation.")
	}

	// Create a method with the options set.
	service, err := builder.NewService("Library").AddMethod(
		builder.NewMethod(
			"WriteBook",
			builder.RpcTypeMessage(builder.NewMessage("WriteBookRequest"), false),
			builder.RpcTypeMessage(builder.NewMessage("WriteBookResponse"), false),
		).SetOptions(opts),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build service.")
	}

	// Establish that we get back both HTTP rules.
	resp := GetHTTPRules(service.GetMethods()[0])
	if got, want := resp[0], httpRule; got != want {
		t.Errorf("Expected the first rule to be %v; got %v.", want, got)
	}
	if got, want := resp[1], httpRule.GetAdditionalBindings()[0]; got != want {
		t.Errorf("Expected the second rule ot be %v; got %v.", want, got)
	}
}

func TestGetHTTPRulesEmpty(t *testing.T) {
	// Create a method with no actual HTTP rules.
	service, err := builder.NewService("Library").AddMethod(
		builder.NewMethod(
			"WriteBook",
			builder.RpcTypeMessage(builder.NewMessage("WriteBookRequest"), false),
			builder.RpcTypeMessage(builder.NewMessage("WriteBookResponse"), false),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build service.")
	}

	// Establish that we get back an empty list of rules.
	if got, want := GetHTTPRules(service.GetMethods()[0]), []*annotations.HttpRule{}; !reflect.DeepEqual(got, want) {
		t.Errorf("Expected empty slice of rules; got %v.", got)
	}
}

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

package aip0122

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func TestHttpUriField(t *testing.T) {
	tests := []struct {
		testName   string
		input      string
		output     string
		httpRule   *annotations.HttpRule
		methodName string
		msg        string
	}{
		{
			"Valid",
			"CreateInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{parent=publishers/*/instanceGroups/*}",
				},
			},
			"CreateInstanceGroup",
			"",
		},
		{
			"InValid-Get",
			"GetInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{
					Get: "/v1/{name=publishers/*/instance_groups/*}",
				},
			},
			"GetInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
		{
			"InValid-Create",
			"CreateInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{parent=publishers/*/instance_groups/*}",
				},
			},
			"CreateInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
		{
			"InValid-Delete",
			"DeleteInstanceGroupsRequest",
			"Empty",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Delete{
					Delete: "/v1/{name=publishers/*/instance_groups/*}",
				},
			},
			"DeleteInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
		{
			"InValid-Update_Patch",
			"UpdateInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Patch{
					Patch: "/v1/{name=publishers/*/instance_groups/*}",
				},
			},
			"UpdateInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
		{
			"InValid-Update_Put",
			"UpdateInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Put{
					Put: "/v1/{name=publishers/*/instance_groups/*}",
				},
			},
			"UpdateInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
		{
			"InValid-Custom",
			"MoveInstanceGroupsRequest",
			"InstanceGroups",
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Custom{
					Custom: &annotations.CustomHttpPattern{
						Path: "/v1/{move=publishers/*/instance_groups/*}",
					},
				},
			},
			"CreateInstanceGroup",
			"HTTP URL pattern shouldn't",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create a MethodOptions with the annotation set.
			opts := &dpb.MethodOptions{}

			if err := proto.SetExtension(opts, annotations.E_Http, test.httpRule); err != nil {
				t.Fatalf("Failed to set google.api.http annotation.")
			}

			// Create a minimal service with a AIP-0122
			service, err := builder.NewService("InstanceGroup").AddMethod(builder.NewMethod(test.methodName,
				builder.RpcTypeMessage(builder.NewMessage(test.input), false),
				builder.RpcTypeMessage(builder.NewMessage(test.output), false),
			).SetOptions(opts)).Build()
			if err != nil {
				t.Fatalf("Could not build %s method.", test.methodName)
			}

			// Run the method, ensure we get what we expect.
			problems := httpURLPattern.Lint(service.GetFile())
			if test.msg == "" && len(problems) > 0 {
				t.Errorf("Got %v, expected no problems.", problems)
			} else if test.msg != "" && len(problems) == 0 {
				t.Errorf("Got no problems, expected 1.")
			} else if test.msg != "" && !strings.Contains(problems[0].Message, test.msg) {
				t.Errorf("Got %q, expected message containing %q", problems[0].Message, test.msg)
			}
		})
	}
}

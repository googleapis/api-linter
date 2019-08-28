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

package rules

// import (
// 	"google.golang.org/protobuf/reflect/protoreflect"
// 	"github.com/googleapis/api-linter/lint"
// 	"github.com/googleapis/api-linter/rules/descriptor"
// )

// func init() {
// 	registerRules(checkProtoVersion())
// }

// // checkProtoVersion returns a lint.Rule
// // that checks if an API is using proto3.
// func checkProtoVersion() lint.Rule {
// 	return &descriptor.CallbackRule{
// 		RuleInfo: lint.RuleInfo{
// 			Name:         lint.NewRuleName("core", "proto_version"),
// 			Description:  "APIs should use proto3",
// 			URI:          `https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/proto-version.md?cl=head`,
// 			RequestTypes: []lint.RequestType{lint.ProtoRequest},
// 		},
// 		Callback: descriptor.Callbacks{
// 			FileCallback: func(f protoreflect.FileDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
// 				location, _ := s.SyntaxLocation()
// 				if f.Syntax() != protoreflect.Proto3 {
// 					return []lint.Problem{
// 						{
// 							Message:    "APIs should use proto3",
// 							Suggestion: "proto3",
// 							Location:   location,
// 						},
// 					}, nil
// 				}
// 				return nil, nil
// 			},
// 		},
// 	}
// }

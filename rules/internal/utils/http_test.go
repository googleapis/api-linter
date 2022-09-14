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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/rules/internal/testutils"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

func TestGetHTTPRules(t *testing.T) {
	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
		t.Run(method, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				service Library {
					rpc FrobBook(FrobBookRequest) returns (FrobBookResponse) {
						option (google.api.http) = {
							{{.M}}: "/v1/publishers/*/books/*"
							additional_bindings {
								{{.M}}: "/v1/books/*"
							}
						};
					}
				}
				message FrobBookRequest {}
				message FrobBookResponse {}
			`, struct{ M string }{M: strings.ToLower(method)})

			// Get the rules.
			resp := GetHTTPRules(file.GetServices()[0].GetMethods()[0])

			// Establish that we get back both HTTP rules, in order.
			if got, want := resp[0].URI, "/v1/publishers/*/books/*"; got != want {
				t.Errorf("Rule 1: Got URI %q, expected %q.", got, want)
			}
			if got, want := resp[1].URI, "/v1/books/*"; got != want {
				t.Errorf("Rule 2: Got URI %q, expected %q.", got, want)
			}
			for _, httpRules := range resp {
				if got, want := httpRules.Method, method; got != want {
					t.Errorf("Got method %q, expected %q.", got, want)
				}
			}
		})
	}
}

func TestGetHTTPRulesEmpty(t *testing.T) {
	file := testutils.ParseProto3String(t, `
		import "google/api/annotations.proto";
		service Library {
			rpc FrobBook(FrobBookRequest) returns (FrobBookResponse);
		}
		message FrobBookRequest {}
		message FrobBookResponse {}
	`)
	if resp := GetHTTPRules(file.GetServices()[0].GetMethods()[0]); len(resp) > 0 {
		t.Errorf("Got %v; expected no rules.", resp)
	}
}

func TestParseRuleEmpty(t *testing.T) {
	http := &apb.HttpRule{}
	if got := parseRule(http); got != nil {
		t.Errorf("Got %v, expected nil.", got)
	}
}

func TestGetHTTPRulesCustom(t *testing.T) {
	file := testutils.ParseProto3String(t, `
		import "google/api/annotations.proto";
		service Library {
			rpc FrobBook(FrobBookRequest) returns (FrobBookResponse) {
				option (google.api.http) = {
					custom: {
						kind: "HEAD"
						path: "/v1/books/*"
					}
				};
			}
		}
		message FrobBookRequest {}
		message FrobBookResponse {}
	`)
	rule := GetHTTPRules(file.GetServices()[0].GetMethods()[0])[0]
	if got, want := rule.Method, "HEAD"; got != want {
		t.Errorf("Got %q; expected %q.", got, want)
	}
	if got, want := rule.URI, "/v1/books/*"; got != want {
		t.Errorf("Got %q; expected %q.", got, want)
	}
}

func TestGetPlainURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		plainURI string
	}{
		{"KeyOnly", "/v1/publishers/{pub_id}/books/{book_id}", "/v1/publishers/*/books/*"},
		{"KeyValue", "/v1/{name=publishers/*/books/*}", "/v1/publishers/*/books/*"},
		{"MultiKeyValue", "/v1/{publisher=publishers/*}/{book=books/*}", "/v1/publishers/*/books/*"},
		{"TemplateVariableSegment", "/{$api_version}/publishers/{pub_id}/books/{book_id}", "/v/publishers/*/books/*"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &HTTPRule{URI: test.uri}
			if diff := cmp.Diff(rule.GetPlainURI(), test.plainURI); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestGetVariables(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		vars map[string]string
	}{
		{"KeyOnly", "/v1/publishers/{pub_id}/books/{book_id}", map[string]string{"pub_id": "*", "book_id": "*"}},
		{"KeyValue", "/v1/{name=publishers/*/books/*}", map[string]string{"name": "publishers/*/books/*"}},
		{"MultiKeyValue", "/v1/{publisher=publishers/*}/{book=books/*}", map[string]string{"publisher": "publishers/*", "book": "books/*"}},
		{"IgnoreVersioningVariable", "/{$api_version}/{name=publishers/*/books/*}", map[string]string{"name": "publishers/*/books/*"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rule := &HTTPRule{URI: test.uri}
			if diff := cmp.Diff(rule.GetVariables(), test.vars); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHasHTTPRules(t *testing.T) {
	for _, tst := range []struct {
		name       string
		Annotation string
	}{
		{"has_rule", `option (google.api.http) = {get: "/v1/foos"};`},
		{"no_rule", ""},
	} {
		t.Run(tst.name, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
					
				service Foo {
					rpc ListFoos (ListFoosRequest) returns (ListFoosResponse) {
						{{ .Annotation }}
					}
				}

				message ListFoosRequest {}
				message ListFoosResponse {}
			`, tst)

			want := tst.Annotation != ""
			if got := HasHTTPRules(file.GetServices()[0].GetMethods()[0]); got != want {
				t.Errorf("Got %v, expected %v", got, want)
			}
		})
	}
}

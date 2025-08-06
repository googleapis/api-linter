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

package aip0123

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
	apb "google.golang.org/genproto/googleapis/api/annotations"
)

func TestAddRules(t *testing.T) {
	if err := AddRules(lint.NewRuleRegistry()); err != nil {
		t.Errorf("AddRules got an error: %v", err)
	}
}

func TestIsNestedName(t *testing.T) {
	for _, test := range []struct {
		name     string
		resource *apb.ResourceDescriptor
		want     bool
	}{
		{
			name: "top level",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Project",
				Pattern:  []string{"projects/{project}"},
				Singular: "project",
				Plural:   "projects",
			},
			want: false,
		},
		{
			name: "top level with extra leading slash",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Project",
				Pattern:  []string{"/projects/{project}"},
				Singular: "project",
				Plural:   "projects",
			},
			want: false,
		},
		{
			name: "top level with extra trailing slash",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Project",
				Pattern:  []string{"projects/{project}/"},
				Singular: "project",
				Plural:   "projects",
			},
			want: false,
		},
		{
			name: "non-nested child collection",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Location",
				Pattern:  []string{"projects/{project}/locations/{location}"},
				Singular: "location",
				Plural:   "locations",
			},
			want: false,
		},
		{
			name: "non-nested child collection multi-word",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/BillingAccount",
				Pattern:  []string{"projects/{project}/billingAccounts/{billing_account}"},
				Singular: "billingAccount",
				Plural:   "billingAccounts",
			},
			want: false,
		},
		{
			name: "nested child collection full",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserEvent",
				Pattern:  []string{"projects/{project}/users/{user}/userEvents/{user_event}"},
				Singular: "userEvent",
				Plural:   "userEvents",
			},
			want: true,
		},
		{
			name: "nested child collection reduced",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserEvent",
				Pattern:  []string{"projects/{project}/users/{user}/events/{event}"},
				Singular: "userEvent",
				Plural:   "userEvents",
			},
			want: true,
		},
		{
			name: "nested singleton full",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserConfig",
				Pattern:  []string{"projects/{project}/users/{user}/userConfig"},
				Singular: "userConfig",
				Plural:   "userConfigs",
			},
			want: true,
		},
		{
			name: "nested singleton reduced",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserConfig",
				Pattern:  []string{"projects/{project}/users/{user}/config"},
				Singular: "userConfig",
				Plural:   "userConfigs",
			},
			want: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := isNestedName(test.resource); got != test.want {
				t.Errorf("got %v, expected %v for pattern %q", got, test.want, test.resource.GetPattern()[0])
			}
		})
	}
}

func TestNestedSingular(t *testing.T) {
	for _, test := range []struct {
		name     string
		resource *apb.ResourceDescriptor
		want     string
	}{
		{
			name: "top level",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Project",
				Pattern:  []string{"projects/{project}"},
				Singular: "project",
				Plural:   "projects",
			},
		},
		{
			name: "non-nested child collection",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/Location",
				Pattern:  []string{"projects/{project}/locations/{location}"},
				Singular: "location",
				Plural:   "locations",
			},
		},
		{
			name: "non-nested child collection multi-word",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/BillingAccount",
				Pattern:  []string{"projects/{project}/billingAccounts/{billing_account}"},
				Singular: "billingAccount",
				Plural:   "billingAccounts",
			},
		},
		{
			name: "nested child collection full",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserEvent",
				Pattern:  []string{"projects/{project}/users/{user}/userEvents/{user_event}"},
				Singular: "userEvent",
				Plural:   "userEvents",
			},
			want: "event",
		},
		{
			name: "nested singleton full",
			resource: &apb.ResourceDescriptor{
				Type:     "example.googleapis.com/UserConfig",
				Pattern:  []string{"projects/{project}/users/{user}/userConfig"},
				Singular: "userConfig",
				Plural:   "userConfigs",
			},
			want: "config",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := nestedSingular(test.resource); got != test.want {
				t.Errorf("got %v, expected %v for pattern %q", got, test.want, test.resource.GetPattern()[0])
			}
		})
	}
}

func TestIsTopLevelResourcePattern(t *testing.T) {
	for _, test := range []struct {
		name    string
		pattern string
		want    bool
	}{
		{
			name:    "top level",
			pattern: "projects/{project}",
			want:    true,
		},
		{
			name:    "not top level",
			pattern: "projects/{project}/locations/{location}",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := isRootLevelResourcePattern(test.pattern); got != test.want {
				t.Errorf("got %v, expected %v for pattern %q", got, test.want, test.pattern)
			}
		})
	}
}

func TestGetParentIDVariable(t *testing.T) {
	for _, test := range []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name:    "top level",
			pattern: "projects/{project}",
		},
		{
			name:    "no variables",
			pattern: "foos",
		},
		{
			name:    "standard child collection",
			pattern: "projects/{project}/locations/{location}",
			want:    "project",
		},
		{
			name:    "singleton",
			pattern: "projects/{project}/config",
			want:    "project",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := getParentIDVariable(test.pattern); got != test.want {
				t.Errorf("got %v, expected %v for pattern %q", got, test.want, test.pattern)
			}
		})
	}
}

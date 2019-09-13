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

// RuleMetadata defines a common structure for name, description, and URI.
type RuleMetadata struct {
	// The name of the rule.
	Name RuleName

	// The human-readable description of the rule.
	Description string

	// The URI where the guidance is documented.
	URI string
}

// GetName returns the rule's name.
func (rm *RuleMetadata) GetName() RuleName {
	return rm.Name
}

// GetDescription returns the rule's description.
func (rm *RuleMetadata) GetDescription() string {
	return rm.Description
}

// GetURI returns the rule's URI.
func (rm *RuleMetadata) GetURI() string {
	return rm.URI
}

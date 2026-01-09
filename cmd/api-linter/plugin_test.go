// Copyright 2025 Google LLC
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

package main

import (
	"testing"

	"github.com/googleapis/api-linter/lint"
)

func TestLoadCustomRulePlugins(t *testing.T) {
	// Mock registry
	mockRegistry := lint.NewRuleRegistry()

	// Test with empty paths
	if err := loadCustomRulePlugins(nil, mockRegistry); err != nil {
		t.Fatalf("loadCustomRulePlugins(nil) = %v; want no error", err)
	}

	// Test with invalid paths - should return an error
	invalidPaths := []string{"nonexistent.so"}
	err := loadCustomRulePlugins(invalidPaths, mockRegistry)
	if err == nil {
		t.Fatalf("loadCustomRulePlugins(%v) = nil; want error", invalidPaths)
	}

	// Note: Testing with actual plugins would require building real .so files,
	// which is complex to set up in unit tests. In a real integration test,
	// we would build a test plugin and verify it loads correctly.
}

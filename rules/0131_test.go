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

package rules

import (
	"testing"

	"github.com/jhump/protoreflect/desc/builder"
)

func TestCheckGetMessageNameValid(t *testing.T) {
	// Create an appropriate method, with a correct message name.
	service, err := builder.NewService("Library").AddMethod(builder.NewMethod("GetBook",
		builder.RpcTypeMessage(builder.NewMessage("GetBookRequest"), false),
		builder.RpcTypeMessage(builder.NewMessage("Book"), false),
	)).Build()
	if err != nil {
		t.Fatalf("Could not build GetBook method.")
	}

	// Run the lint rule; it should return no problems.
	if problems := checkGetRequestMessageName.LintMethod(service.GetMethods()[0]); len(problems) > 0 {
		t.Errorf("False positive on rule %s: %#v", checkGetRequestMessageName.Name, problems)
	}
}

func TestCheckGetMessageNameInvalid(t *testing.T) {
	// Create an appropriate method, with an incorrect request message name.
	service, err := builder.NewService("Library").AddMethod(builder.NewMethod("GetBook",
		builder.RpcTypeMessage(builder.NewMessage("Book"), false),
		builder.RpcTypeMessage(builder.NewMessage("Book"), false),
	)).Build()
	if err != nil {
		t.Fatalf("Could not build GetBook method.")
	}

	// Run the lint rule; it should return no problems.
	if problems := checkGetRequestMessageName.LintMethod(service.GetMethods()[0]); len(problems) < 1 {
		t.Errorf("False negative on rule %s: %#v", checkGetRequestMessageName.Name, problems)
	}
}

func TestCheckGetMessageNameIrrelevant(t *testing.T) {
	// Create an appropriate method, with a correct message name.
	service, err := builder.NewService("Library").AddMethod(builder.NewMethod("AcquireBook",
		builder.RpcTypeMessage(builder.NewMessage("AcquireBookReq"), false),
		builder.RpcTypeMessage(builder.NewMessage("Book"), false),
	)).Build()
	if err != nil {
		t.Fatalf("Could not build AcquireBook method.")
	}

	// Run the lint rule; it should return no problems.
	if problems := checkGetRequestMessageName.LintMethod(service.GetMethods()[0]); len(problems) > 0 {
		t.Errorf(
			"False positive on rule %s (should have been irrelevant): %#v",
			checkGetRequestMessageName.Name,
			problems,
		)
	}
}

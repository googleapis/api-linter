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

// Package aip0192 contains rules defined in https://aip.dev/192.
package aip0192

import (
	"strings"

	"github.com/googleapis/api-linter/lint"
)

// AddRules adds all of the AIP-192 rules to the provided registry.
func AddRules(r lint.RuleRegistry) {
	r.Register(
		onlyLeadingComments,
	)
}

func separateInternalComments(comments ...string) struct {
	Internal []string
	External []string
} {
	answer := struct {
		Internal []string
		External []string
	}{}
	for _, c := range comments {
		for len(c) > 0 {
			// Anything before the `(--` is external string.
			open := strings.SplitN(c, "(--", 2)
			if ex := strings.TrimSpace(open[0]); ex != "" {
				answer.External = append(answer.External, ex)
			}
			if len(open) > 1 {
				c = strings.TrimSpace(open[1])
			} else {
				c = ""
			}

			// Now that the opening component is tokenized, anything before
			// the `--)` is internal string.
			close := strings.SplitN(c, "--)", 2)
			if in := strings.TrimSpace(close[0]); in != "" {
				answer.Internal = append(answer.Internal, in)
			}
			if len(close) > 1 {
				c = strings.TrimSpace(close[1])
			} else {
				c = ""
			}
		}
	}
	return answer
}

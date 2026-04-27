// Copyright 2022 Google LLC
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
	"bytes"
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/v2/lint"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// formatGitHubActionOutput returns lint errors in GitHub actions format.
func formatGitHubActionOutput(responses []lint.Response) []byte {
	var buf bytes.Buffer
	for _, response := range responses {
		for _, problem := range response.Problems {
			// lint example:
			// ::error file={name},line={line},endLine={endLine},title={title}::{message}
			// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message

			fmt.Fprintf(&buf, "::error file=%s", response.FilePath)

			if problem.Location != nil {
				start, end := fileLocationFromPBLocation(problem.Location)
				fmt.Fprintf(&buf, ",line=%d,endLine=%d,col=%d,endColumn=%d", start.Line, end.Line, start.Column, end.Column)
			}

			// GitHub uses :: as control characters (which are also used to delimit
			// linter rules. In order to prevent confusion, replace the double colon
			// with two Armenian full stops which are indistinguishable to my eye.
			runeThatLooksLikeTwoColonsButIsActuallyTwoArmenianFullStops := "։։"
			title := strings.ReplaceAll(string(problem.RuleID), "::", runeThatLooksLikeTwoColonsButIsActuallyTwoArmenianFullStops)
			message := strings.ReplaceAll(problem.Message, "\n", "%0A")
			uri := problem.GetRuleURI()
			if uri != "" {
				message += "%0A%0A" + uri
			}
			fmt.Fprintf(&buf, ",title=%s::%s\n", title, message)
		}
	}

	return buf.Bytes()
}

type position struct {
	Line   int
	Column int
}

// Implementation copied from lint/problem.go
func fileLocationFromPBLocation(l *dpb.SourceCodeInfo_Location) (start, end position) {
	start = position{
		Line:   int(l.Span[0]) + 1,
		Column: int(l.Span[1]) + 1,
	}

	if len(l.Span) == 4 {
		end = position{
			Line:   int(l.Span[2]) + 1,
			Column: int(l.Span[3]),
		}
	} else {
		end = position{
			Line:   int(l.Span[0]) + 1,
			Column: int(l.Span[2]),
		}
	}
	return start, end
}

// Copyright 2020 Google LLC
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

package aip0192

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules/internal/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var trademarkAliases = map[string][]string{
	"App Engine":     {"GAE", "gae", "AppEngine", "App engine"},
	"BigQuery":       {"Bigquery", "Big Query", "BQ"},
	"BigQuery ML":    {"BQML", "bqml"},
	"Bigtable":       {"BigTable", "Big Table", "Big table"},
	"Bitbucket":      {"BitBucket", "Bit Bucket"},
	"Cloud Storage":  {"GCS", "gcs"},
	"Compute Engine": {"GCE", "gce"},
	"Dataflow":       {"Data Flow", "Data flow", "DataFlow"},
	"Dataprep":       {"Data Prep", "Data prep", "DataPrep"},
	"Dialogflow":     {"DialogFlow", "Dialog Flow", "Dialog flow"},
	"Directory Sync": {"GCDS", "CDS", "gcds", "cds", "DirectorySync"},
	"GitHub":         {"Github", "Git Hub"},
	"GitLab":         {"Gitlab", "Git Lab"},
	"G Suite":        {"GSuite", "G-Suite", "gSuite"},
	"Pub/Sub":        {"PubSub", "Pubsub", "Cloud Pub/Sub"},
	"Service Mesh":   {"ASM", "CSM", "GCSM", "csm", "gcsm"},
	"Stack Overflow": {"StackOverflow"},
}

// We actually want regexes so we do not accidentally false-positive acronyms
// that *contain* our matches. (For example, "BQD" should not match and tell us
// to change to BigQuery.)
func defaultTrademarkTypos() map[string][]*regexp.Regexp {
	tmRegexes := map[string][]*regexp.Regexp{}
	for k, tms := range trademarkAliases {
		tmReg := []*regexp.Regexp{}
		for _, tm := range tms {
			tmReg = append(tmReg, regexp.MustCompile(`\b`+strings.ReplaceAll(tm, " ", `\s+`)+`\b`))
		}
		tmRegexes[k] = tmReg
	}
	return tmRegexes
}

var tmRegexes = defaultTrademarkTypos()

var trademarkedNames = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "trademarked-names"),
	LintDescriptor: func(d protoreflect.Descriptor) (problems []lint.Problem) {
		c := strings.Join(
			utils.SeparateInternalComments(d.GetSourceInfo().GetLeadingComments()).External,
			"\n",
		)
		for want, badThings := range tmRegexes {
			for _, bad := range badThings {
				if bad.MatchString(c) {
					problems = append(problems, lint.Problem{
						Message:    fmt.Sprintf("Use %q in comments, not %q.", want, bad),
						Descriptor: d,
					})
				}
			}
		}
		return
	},
}

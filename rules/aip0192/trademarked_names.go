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
	"github.com/jhump/protoreflect/desc"
)

var trademarkAliases = map[string][]string{
	"App Engine":     []string{"GAE", "gae", "AppEngine", "App engine"},
	"BigQuery":       []string{"Bigquery", "Big Query", "BQ"},
	"BigQuery ML":    []string{"BQML", "bqml"},
	"Bigtable":       []string{"BigTable", "Big Table", "Big table"},
	"Bitbucket":      []string{"BitBucket", "Bit Bucket"},
	"Cloud Storage":  []string{"GCS", "gcs"},
	"Compute Engine": []string{"GCE", "gce"},
	"Dataflow":       []string{"Data Flow", "Data flow", "DataFlow"},
	"Dataprep":       []string{"Data Prep", "Data prep", "DataPrep"},
	"Dialogflow":     []string{"DialogFlow", "Dialog Flow", "Dialog flow"},
	"Directory Sync": []string{"GCDS", "CDS", "gcds", "cds", "DirectorySync"},
	"GitHub":         []string{"Github", "Git Hub"},
	"GitLab":         []string{"Gitlab", "Git Lab"},
	"G Suite":        []string{"GSuite", "G-Suite", "gSuite"},
	"Pub/Sub":        []string{"PubSub", "Pubsub", "Cloud Pub/Sub"},
	"Service Mesh":   []string{"ASM", "CSM", "GCSM", "csm", "gcsm"},
	"Stack Overflow": []string{"StackOverflow"},
}

// We actually want regexes so we do not accidentally false-positive acronyms
// that *contain* our matches. (For example, "BQD" should not match and tell us
// to change to BigQuery.)
var tmRegexes = map[string][]*regexp.Regexp{}

func defaultTrademarkTypos() {
	for k, tms := range trademarkAliases {
		tmReg := []*regexp.Regexp{}
		for _, tm := range tms {
			tmReg = append(tmReg, regexp.MustCompile(`\b`+strings.ReplaceAll(tm, " ", `\s+`)+`\b`))
		}
		tmRegexes[k] = tmReg
	}
}

func init() {
	defaultTrademarkTypos()
}

var trademarkedNames = &lint.DescriptorRule{
	Name: lint.NewRuleName(192, "trademarked-names"),
	LintDescriptor: func(d desc.Descriptor) (problems []lint.Problem) {
		c := strings.Join(
			separateInternalComments(d.GetSourceInfo().GetLeadingComments()).External,
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

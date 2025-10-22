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

import (
	"encoding/json"

	"google.golang.org/protobuf/reflect/protoreflect"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

// Problem contains information about a result produced by an API Linter.
//
// All rules return []Problem. Most lint rules return 0 or 1 problems, but
// occasionally there are rules that may return more than one.
type Problem struct {
	// Message provides a short description of the problem.
	// This should be no more than a single sentence.
	Message string

	// Suggestion provides a suggested fix, if applicable.
	//
	// This integrates with certain IDEs to provide "push-button" fixes,
	// so these need to be machine-readable, not just human-readable.
	// Additionally, when setting `Suggestion`, one should almost always set
	// `Location` also, to ensure that the text being replaced is sufficiently
	// precise.
	Suggestion string

	// Descriptor provides the descriptor related to the problem. This must be
	// set on every Problem.
	//
	// If `Location` is not specified, then the starting location of
	// the descriptor is used as the location of the problem.
	Descriptor protoreflect.Descriptor

	// Location provides the location of the problem.
	//
	// If unset, the location of the descriptor is used.
	// This should almost always be set if `Suggestion` is set. The best way to
	// do this is by using the helper methods in `location.go`.
	Location *dpb.SourceCodeInfo_Location

	// RuleID provides the ID of the rule that this problem belongs to.
	// DO NOT SET: The linter sets this automatically.
	RuleID RuleName // FIXME: Make this private (cmd/summary_cli.go is the challenge).

	// The category for this problem, based on user configuration.
	category string

	//lint:ignore U1000 ignored via golint previously
	noPositional struct{}
}

// MarshalJSON defines how to represent a Problem in JSON.
func (p Problem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.marshal())
}

// MarshalYAML defines how to represent a Problem in YAML.
func (p Problem) MarshalYAML() (interface{}, error) {
	return p.marshal(), nil
}

// Marshal defines how to represent a serialized Problem.
func (p Problem) marshal() interface{} {
	var fl fileLocation
	if p.Location != nil {
		// If Location is set, use it.
		fl = fileLocationFromPBLocation(p.Location, p.Descriptor)
	} else if p.Descriptor != nil {
		// Otherwise, use the descriptor's location.
		// This is the protobuf-go idiomatic way to get the source location.
		// Note: ParentFile() called on a FileDescriptor returns itself.
		loc := p.Descriptor.ParentFile().SourceLocations().ByDescriptor(p.Descriptor)
		fl = fileLocation{
			Path: p.Descriptor.ParentFile().Path(),
			Start: position{
				Line:   loc.StartLine + 1,
				Column: loc.StartColumn + 1,
			},
			End: position{
				Line:   loc.EndLine + 1,
				Column: loc.EndColumn,
			},
		}
	} else {
		// Default location if no descriptor.
		fl = fileLocationFromPBLocation(nil, nil)
	}

	// Return a marshal-able structure.
	return struct {
		Message    string       `json:"message" yaml:"message"`
		Suggestion string       `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
		Location   fileLocation `json:"location" yaml:"location"`
		RuleID     RuleName     `json:"rule_id" yaml:"rule_id"`
		RuleDocURI string       `json:"rule_doc_uri" yaml:"rule_doc_uri"`
		Category   string       `json:"category,omitempty" yaml:"category,omitempty"`
	}{
		p.Message,
		p.Suggestion,
		fl,
		p.RuleID,
		p.GetRuleURI(),
		p.category,
	}
}

// GetRuleURI returns a URI to learn more about the problem.
func (p Problem) GetRuleURI() string {
	return getRuleURL(string(p.RuleID), ruleURLMappings)
}

// position describes a one-based position in a source code file.
// They are one-indexed, as a human counts lines or columns.
type position struct {
	Line   int `json:"line_number" yaml:"line_number"`
	Column int `json:"column_number" yaml:"column_number"`
}

// fileLocation describes a location in a source code file.
//
// Note: Positions are one-indexed, as a human counts lines or columns
// in a file.
type fileLocation struct {
	Start position `json:"start_position" yaml:"start_position"`
	End   position `json:"end_position" yaml:"end_position"`
	Path  string   `json:"path" yaml:"path"`
}

// fileLocationFromPBLocation returns a new fileLocation object based on a
// protocol buffer SourceCodeInfo_Location
func fileLocationFromPBLocation(l *dpb.SourceCodeInfo_Location, d protoreflect.Descriptor) fileLocation {
	// Spans are guaranteed by protobuf to have either three or four ints.
	span := []int32{0, 0, 1}
	if l != nil {
		span = l.Span
	}

	var fl fileLocation
	if d != nil {
		fl = fileLocation{Path: d.ParentFile().Path()}
	}

	// If `span` has four ints; they correspond to
	// [start line, start column, end line, end column].
	//
	// We add one because spans are zero-indexed, but not to the end column
	// because we want the ending position to be inclusive and not exclusive.
	if len(span) == 4 {
		fl.Start = position{
			Line:   int(span[0]) + 1,
			Column: int(span[1]) + 1,
		}
		fl.End = position{
			Line:   int(span[2]) + 1,
			Column: int(span[3]),
		}
		return fl
	}

	// Okay, `span` has three ints; they correspond to
	// [start line, start column, end column].
	//
	// We add one because spans are zero-indexed, but not to the end column
	// because we want the ending position to be inclusive and not exclusive.
	fl.Start = position{
		Line:   int(span[0]) + 1,
		Column: int(span[1]) + 1,
	}
	fl.End = position{
		Line:   int(span[0]) + 1,
		Column: int(span[2]),
	}
	return fl
}

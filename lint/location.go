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

import "fmt"

// Location describes a location in a source code file.
//
// Note: positions are one-based.
type Location struct {
	Start Position `json:"start_position" yaml:"start_position"`
	End   Position `json:"end_position" yaml:"end_position"`
}

// IsValid checks if the location is constructed properly and
// returns true if so.
func (l Location) IsValid() bool {
	return l.Start.IsValid() &&
		l.End.IsValid() &&
		(l.End.Line > l.Start.Line ||
			l.End.Line == l.Start.Line && l.End.Column >= l.Start.Column)
}

// String returns the string representation.
func (l Location) String() string {
	return fmt.Sprintf("{start: %s, end: %s}", l.Start, l.End)
}

// Position describes a one-based position in a source code file.
type Position struct {
	Line   int `json:"line_number" yaml:"line_number"`
	Column int `json:"column_number" yaml:"column_number"`
}

// IsValid checks if the position is constructed properly and
// returns true if so.
func (p Position) IsValid() bool {
	return p.Line >= 1 && p.Column >= 1
}

// String returns the string representation.
func (p Position) String() string {
	return fmt.Sprintf("{line: %d, column: %d}", p.Line, p.Column)
}

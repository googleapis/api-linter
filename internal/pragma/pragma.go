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

// Package pragma provides types that can be embedded into a struct to statically enforce or prevent certain language properties.
package pragma

// DoNotImplement can be embedded in an interface to prevent trivial
// implementations of the interface.
//
// This is useful to prevent unauthorized implementations of an interface
// so that it can be extended in the future for any API changes.
type DoNotImplement interface{ LinterInternal(DoNotImplement) }

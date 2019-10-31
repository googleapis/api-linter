# Google API Linter

The API linter provides real-time checks for compliance with many of Google's
API standards, documented using [API Improvement Proposals][]. It operates on
API surfaces defined in [protocol buffers][].

It identifies common mistakes and inconsistencies in API surfaces:

```proto
// Incorrect.
message GetBookRequest {
  // This is wrong; it should be spelled `name`.
  string book = 1;
}
```

When able, it also offers a suggestion for the correct fix.

**Note:** Not every piece of AIP guidance is able to be expressed as lint rules
(and some things that are able to be expressed may not be written yet). The
linter should be used as a useful tool, but not as a substitute for reading and
understanding API guidance.

Each linter rule has its own [rule documentation][], and rules can be
[configured][configuration] using a config file, or in a proto file itself.

## Installation

To install `api-linter`, use `go get`:

```sh
go get -u github.com/googleapis/api-linter/cmd/api-linter
```

It will install `api-linter` into your local Go binary directory `$HOME/go/bin`.
Ensure that your operating system's `PATH` contains the Go binary directory.

## Usage

```sh
api-linter proto_file1 proto_file2 ...
```

To see the help message, run `api-linter -h`

```text
Usage of api-linter:
  -config string
      The linter config file.
  -output_format string
      The format of the linting results.
      Supported formats include YAML, JSON and summary text.
      YAML is the default.
  -output_path string
      The output file path.
      If not given, the linting results will be printed out to STDOUT.
  -proto_descriptor_set string
      The file descriptor set for searching proto imports.
  -proto_path value
      The folder for searching proto imports.
      May be specified multiple times; directories will be searched in order.
      The current working directory is always used.
```

## Configuration

The linter contains a list of core rules, and by default, they are all enabled.
However, one can disable a rule by using a configuration file or the file
comments.

### Disable a rule using a configuration file

Examples:

* Disable the rule `core::0140::lower-snake` for any proto files under the
directory `tests` using a JSON config file:

```json
[
   {
      "included_paths": ["tests/*.proto"],
      "rule_configs": {
         "core::0140::lower-snake": {"status": "disabled"}
      }
   }
]
```

* Disable the same rule using a YAML config file:

```yaml
---
- included_paths:
  - "**/*.proto"
  rule_configs:
    core::proto_version:
      status: disabled
```

### Disable a rule in the file comments

Examples:

* Disable the rule `core::0140::lower-snake` for the entire file:

```protobuf
// The file comments.
// api-linter: core::0140::lower-snake=disabled
// The above comment will disable the rule
// core::0140::lower-snake for the entire file.

syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    string anotherBadFieldName = 2;
}
```

* Disable the same rule only for a field in its leading comments:

```protobuf
syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    // The leading comments for the field `anotherBadFieldName`
    // api-linter: core::0140::lower-snake=disabled
    string anotherBadFieldName = 2;
}
```

## License

This software is made available under the [Apache 2.0][] license.

[apache 2.0]: https://www.apache.org/licenses/LICENSE-2.0
[api improvement proposals]: https://aip.dev/
[protocol buffers]: https://developers.google.com/protocol-buffers
[rule documentation]: ./rules/index.md

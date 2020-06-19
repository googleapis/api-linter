---
---

# Google API Linter

![ci](https://github.com/googleapis/api-linter/workflows/ci/badge.svg)
![latest release](https://img.shields.io/github/v/release/googleapis/api-linter)
![go version](https://img.shields.io/github/go-mod/go-version/googleapis/api-linter)

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

It will install `api-linter` into your local Go binary directory
`$HOME/go/bin`. Ensure that your operating system's `PATH` contains the Go
binary directory.

**Note:** For working in Google-internal source control, you should use the
released binary `/google/bin/releases/api-linter/api-linter`.

## Usage

```sh
api-linter proto_file1 proto_file2 ...
```

To see the help message, run `api-linter -h`

```text
Usage of api-linter:
      --config string                   The linter config file.
      --descriptor-set-in stringArray   The file containing a FileDescriptorSet for searching proto imports.
                                        May be specified multiple times.
      --disable-rule stringArray        Disable a rule with the given name.
                                        May be specified multiple times.
      --enable-rule stringArray         Enable a rule with the given name.
                                        May be specified multiple times.
      --output-format string            The format of the linting results.
                                        Supported formats include "yaml", "json" and "summary" table.
                                        YAML is the default.
  -o, --output-path string              The output file path.
                                        If not given, the linting results will be printed out to STDOUT.
  -I, --proto-path stringArray          The folder for searching proto imports.
                                        May be specified multiple times; directories will be searched in order.
                                        The current working directory is always used.
```

## License

This software is made available under the [Apache 2.0][] license.

[apache 2.0]: https://www.apache.org/licenses/LICENSE-2.0
[api improvement proposals]: https://aip.dev/
[configuration]: ./configuration.md
[protocol buffers]: https://developers.google.com/protocol-buffers
[rule documentation]: ./rules/index.md

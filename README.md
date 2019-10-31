# API Linter

API linter checks APIs defined in protobuf files. It follows [Google API Design Guide](https://cloud.google.com/apis/design/).

## Requirements

* Install `git` from [https://git-scm.com](https://git-scm.com/);
* Install `go` from [https://golang.org/doc/install](https://golang.org/doc/install);

## Installation

* Install `api-linter` using `go get`:

```sh
go get -u github.com/googleapis/api-linter/cmd/api-linter
```

This installs `api-linter` into your local Go binary folder `$HOME/go/bin`. Ensure that your operating system's `PATH` contains the folder.

## Usage

```sh
api-linter proto_file1 proto_file2 ...
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


## Contributing

To contribute your rules, please open an issue first and follow [those existing rules](https://github.com/googleapis/api-linter/tree/master/rules) as examples. Pull requests are always welcome.

## License

[Apache License 2.0](LICENSE)

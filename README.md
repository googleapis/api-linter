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

Run `api-linter help` to see the usage. Or run `api-linter help checkproto` to see how to check API protobuf files:

```sh
NAME:
   api-linter checkproto - Check protobuf files that define an API

USAGE:
   api-linter checkproto [command options] files...

OPTIONS:
   --cfg value          configuration file path
   --out value          output file path (default: stdout)
   --fmt value          output format (default: "yaml")
   --protoc value       protocol compiler path (default: "protoc")
   --proto_path value   the directory in which for protoc to search for imports (default: ".")
```

See this [example](cmd/api-linter/examples/example.sh).

## Rule Configuration

The linter contains a list of [core rules](rules), and by default, they are all enabled. However, one can disable a rule by using a configuration file or in-file(line) comments.

### Disable a rule using a configuration file

Example:

Disable rule `core::proto_version` for any `.proto` files.

```json
[
   {
      "included_paths": ["**/*.proto"],
      "rule_configs": {
         "core::proto_version": {"status": "disabled"}
      }
   }
]
```

### Disable a rule using in-file(line) comments

Example:

* Disable rule `core::naming_formats::field_names` entirely for a file in the file comments.

```protobuf
// file comments
// api-linter: core::naming_formats::field_names=disabled

syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    string anotherBadFieldName = 2;
}
```

* Disable rule `core::naming_formats::field_names` only for a field in its leading or trailing comments.

```protobuf
syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    // leading comments for field `anotherBadFieldName`
    // api-linter: core::naming_formats::field_names=disabled
    string anotherBadFieldName = 2; // trailing comments (-- api-linter: core::naming_formats::field_names=disabled --)
}
```

## Contributing

To contribute your rules, please open an issue first and follow [those existing rules](https://github.com/googleapis/api-linter/tree/master/rules) as examples. Pull requests are always welcome.

## License

[Apache License 2.0](LICENSE)

# API Linter
API linter checks APIs defined in protobuf files. It follows [Google API Design Guide](https://cloud.google.com/apis/design/).

## Requirements
* Install `git` from [https://git-scm.com](https://git-scm.com/);
* Install `go` from [https://golang.org/doc/install](https://golang.org/doc/install);
* Install `protoc` by following this [guide](http://google.github.io/proto-lens/installing-protoc.html);

## Installation
* Install `api-linter` using `go get`:
```sh
go get -u github.com/googleapis/api-linter/cmd/api-linter
```
* Update the `PATH` environment to include `$HOME/go/bin`.

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

1. Disable rule `core::naming_formats::field_names` entirely for a file in the file comments.
```
// file comments
// (-- api-linter: core::naming_formats::field_names=disabled --)

syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    string anotherBadFieldName = 2;
}
```

2. Disable rule `core::naming_formats::field_names` only for a field in its leading or trailing comments.
```
syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    // leading comments for field `anotherBadFieldName`
    // (-- api-linter: core::naming_formats::field_names=disabled --)
    string anotherBadFieldName = 2; // trailing comments (-- api-linter: core::naming_formats::field_names=disabled --)
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](LICENSE)

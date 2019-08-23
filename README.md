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
// (-- api-linter: core::naming_formats::field_names=disabled --)

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
    // (-- api-linter: core::naming_formats::field_names=disabled --)
    string anotherBadFieldName = 2; // trailing comments (-- api-linter: core::naming_formats::field_names=disabled --)
}
```

## Writing Rules

If you would like to extend the linter with your own logic for linting proto files, you need to
implement the [`Rule` interface][rule_interface]. First, expose the metadata for your rule by
implementing an `Info()` method, which returns a [`RuleInfo` struct][rule_info]. You must give your
rule a unique name, description, and URL, and define the types of requests that the rule should
receive (currently, the only type of request is a [`lint.ProtoRequest`][proto_request_type]).

Then, implement the `Lint` method, which takes a [`lint.Request`][lint_request], and returns a slice
of [`Problem`s][problems] (or an error). The `lint.Request` provides to you a compiled
[`protoreflect.FileDescriptor`][file_descriptor], which allows you to programmatically navigate the
definitions in a proto file, as well as a `lint.DescriptorSource`, which provides some helper
[methods][desc_source] for finding the locations of a given descriptor in the proto source. Please
ensure that every `Problem` you output contains a valid [`lint.Location`][location] (usually
determined through `lint.DescriptorSource`). The `Location` will enable other tools to display
warnings in the relevant location in source code. Furthermore, you can provide a `Suggestion` to
replace the violating code covered by the `Location`.

For Example:

```go
type myRule struct{}

func Info() lint.RuleInfo {
	return lint.RuleInfo{
    Name:         "my_test_rule",
    Description:  "dummy rule that doesn't do anything",
    RequestTypes: []lint.RequestType{lint.ProtoRequest},
    URI:          "http://example.com",
  }
}

func Lint(req lint.Request) ([]lint.Problem, error) {
	if req.ProtoFile().Syntax() != protoreflect.Proto3 {
		loc, err := req.DescriptorSource().SyntaxLocation()
		if err != nil {
			// use start location of file
			loc = lint.Location{
				Start: lint.Position{Line: 1, Column: 1},
				End:   lint.Position{Line: 1, Column: 1},
			}
		}
		return []lint.Problem{{
			Message:    "Please use proto3 syntax.",
			Location:   loc,
			Suggestion: "proto3",
		}}
	}
}
```

[file_descriptor]: https://google.golang.org/protobuf/reflect/protoreflect
[rule_interface]: https://github.com/googleapis/api-linter/blob/master/lint/rule.go
[rule_info]: https://github.com/googleapis/api-linter/blob/master/lint/rule_info.go
[lint_request]: https://github.com/googleapis/api-linter/blob/master/lint/request.go#L25-L28
[problems]: https://github.com/googleapis/api-linter/blob/master/lint/response.go
[proto_request_type]: https://github.com/googleapis/api-linter/blob/master/lint/rule_info.go#L37
[desc_source]: https://github.com/googleapis/api-linter/blob/master/lint/source.go#L151-L174
[location]: https://github.com/googleapis/api-linter/blob/master/lint/location.go

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[Apache License 2.0](LICENSE)

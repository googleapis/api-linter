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
* Update the `$PATH` environment to include `$HOME/go/bin`.

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
   --import_path value  protoc import path (default: ".")
```

See this [example](cmd/api-linter/examples/example.sh).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](LICENSE)

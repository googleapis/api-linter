# API Linter

The API linter is intended to enforce conventions on APIs that are defined with `.proto` files.
The `lint` package implements the core functionality for ingesting a lint request, enumerating
source code, and emitting problems. Any integrations with external error reporting systems will
use this package.

The `protohelpers` package implements some helpers to make linting `.proto` files easier. Although
it's not required, `Rule` implementations will probably want to use this package when they are
linting API definitions.

The `rules` package contains the implementation of specific `Rule`s that are broadly applicable to
APIs defined at Google. New `Rule` implementations should only go here if they are general enough
to be useful across the organization.

## Installation

1. [Install Go](https://golang.org/doc/install)
2. Clone this repository to a location of your choice.
   - If you have your workstation's SSH key linked to your GitHub account:

    ```
    $ git clone git@github.com:googleapis/api-linter.git
    ```

   - Otherwise:
    
    ```
    $ git clone https://github.com/googleapis/api-linter.git
    ```
    
## Running the linter

TODO once we have a `main` CLI package.

## Implementing a Rule

In order to implement a new rule, you must satisfy the [`Rule` interface][rule_interface]. First,
expose the metadata for your rule by implementing an `Info()` method which returns a
[`RuleInfo` struct][rule_info]. You must give your rule a unique name, description, and URL, and
define the types of files that the rule needs to see.

Then, implement the `Lint` method, which takes a [`Request`][lint_request] and returns a
[`Response`][lint_response] (or an error). The `Request` input contains methods that allow your
implementation to inspect a `.proto` file (`ProtoFile()`) and find the location of different proto
constructs in the source code (`DescriptorSource()`). The `Response` output allows you to emit
`Problem`s, which contain a `Message` to the user, a `Location` (in source) of the offending code,
a `Suggestion` of code to replace the offending code with, and the `Descriptor` that caused the
problem.

For example:

```go
type enforceProto3 struct {}

func (r *enforceProto3) Info() lint.RuleInfo {
	return lint.RuleInfo{
		Name:        "enforce_proto_3",
		Description: `This rule enforces that all protofiles use syntax "proto3"`,
		Url:         `http://example.com/foo/bar/`,
		FileTypes:   []FileType{lint.ProtoFile},
	}
}

func (r *enforceProto3) Lint(req lint.Request, source lint.DescriptorSource) (lint.Response, error) {
	if req.ProtoFile().Syntax() != protoreflect.Proto3 {
		location, err := source.SyntaxLocation()
		
		if err != nil {
			return lint.Response{}, err
		}
		
	  return lint.Response{
	  	Problems: []Problem{
	  		{
	  			Message:    "Google APIs should use proto3",
	  			Suggestion: "proto3",
	  			Location:   location,
	  		},
	  	},
	  }
	}
	
	return lint.Response{}, nil
}
```

### Helpers

TODO once location of helpers is finalized.

<!--
The [`protohelpers` package][proto_helpers] provides some abstractions that can make linting proto
files easier in some cases.

#### WalkDescriptor

The `protohelpers.WalkDescriptor(protoreflect.Descriptor, protohelpers.DescriptorConsumer)` function
allows you to pass a `DescriptorConsumer`, which receives each `protoreflect.Descriptor`, one at a
time. If you have particular rules that you want to enforce, but don't want to deal with the logic
of traversing a proto file, `WalkDescriptor` will recursively traverse nested messages, enums,
fields, and any other `Descriptor` found in a proto file, and pass them to your `ConsumeDescriptor`
method. You can use a [type switch][type_switch] to determine the type of descriptor being passed,
and perform whatever other logic needed to implement your rule.

#### DescriptorCallbacks

The [`DescriptorCallbacks` struct][descriptor_callbacks] allows you to implement callbacks for
specific types of descriptors. For example if you want to write a `Rule` that only cares about the
fields defined in a proto, you can write a function that will only receive `FieldDescriptor`s. For
example, the variable `rule` here satisfies the `Rule` interface:

```go
rule := protohelpers.DescriptorCallbacks{
  RuleInfo: lint.RuleInfo{
    Name:        "check_naming_formats.field",
    Description: "check that field names use lower snake case",
    Url:         "https://g3doc.corp.google.com/google/api/tools/linter/g3doc/rules/naming-format.md?cl=head",
    FileTypes:   []lint.FileType{lint.ProtoFile},
  },
  FieldDescriptorCallback: func(d protoreflect.FieldDescriptor, s lint.DescriptorSource) ([]lint.Problem, error) {
    return checkNameFormat(d), nil
  },
}
```
-->


[rule_interface]: https://github.com/googleapis/api-linter/blob/master/lint/rule.go 
[rule_info]: https://github.com/googleapis/api-linter/blob/master/lint/rule_info.go
[lint_request]: https://github.com/googleapis/api-linter/blob/master/lint/request.go
[lint_response]: https://github.com/googleapis/api-linter/blob/master/lint/response.go
[proto_helpers]: https://github.com/googleapis/api-linter/tree/master/protohelpers
[type_switch]: https://tour.golang.org/methods/16
[descriptor_callbacks]: https://github.com/googleapis/api-linter/blob/master/protohelpers/descriptor_callbacks.go
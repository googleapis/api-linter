# Configuration

The linter contains a list of core rules, and by default, they are all enabled.
However, one can disable a rule by using a CLI flag, a configuration file or the file
comments.

## CLI flag

We can disable a rule using the `-disable_rule` flag.

Example:

Disable the rule `core::0140::lower-snake`:

```sh
api-linter -disable_rule core::0140::lower-snake test.proto
```

To disable multiple rules, specify the flag multiple times respectively:

```sh
api-linter -disable_rule core::0140::lower-snake -disable_rule core::0131::request-name-field test.proto
```

## Configuration file

The linter accepts a configuration file using the `-config` CLI switch.

Examples:

Disable the rule `core::0140::lower-snake` for any proto files under the
directory `tests` using a JSON config file:

```json
[
  {
    "included_paths": ["tests/**/*.proto"],
    "disabled_rules": ["core::0140::lower-snake"],
    "import_paths": ["thirty_part"]
  }
]
```

Disable the same rule using a YAML config file:

```yaml
---
- included_paths:
    - 'tests/**/*.proto'
  disabled_rules:
    - 'core::0140::lower-snake'
  import_paths:
    - 'thirty_part'
```

## Proto comments

Examples:

Disable the rule `core::0140::lower-snake` for the entire file:

```proto
// A file comment:
// (-- api-linter: core::0140::lower-snake=disabled --)
//
// The above comment will disable the rule
// `core::0140::lower-snake` for the entire file.

syntax = "proto3";

package google.api.linter.examples;

message Example {
    string badFieldName = 1;
    string anotherBadFieldName = 2;
}
```

Disable the same rule only for a field by using a leading comment:

```protobuf
syntax = "proto3";

package google.api.linter.examples;

message Example {
    // This field will trigger a lint error.
    string badFieldName = 1;

    // This field will not trigger a lint error.
    // (-- api-linter: core::0140::lower-snake=disabled --)
    string anotherBadFieldName = 2;
}
```

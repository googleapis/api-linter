---
rule:
  aip: 191
  name: [core, '0191', java-package]
  summary: All proto files must set `option java_package`.
permalink: /191/java-package
redirect_from:
  - /0191/java-package
---

# Java package annotation

This rule enforces that every proto file for a public API surface sets
`option java_package`, as mandated in [AIP-191][].

## Details

This rule looks at each proto file, and complains if the `java_package` file
annotation is not present, or set to empty string.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

// Needs `option java_package = "com.google.example.v1";`.
option java_multiple_files = true;
option java_outer_classname = "LibraryProto";

message Book {}
```

**Correct** code for this rule:

```proto
// Correct.
syntax = "proto3";

package google.example.v1;

option java_package = "com.google.example.v1";
option java_multiple_files = true;
option java_outer_classname = "LibraryProto";

message Book {}
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::java-package=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

package google.example.v1;

option java_multiple_files = true;
option java_outer_classname = "LibraryProto";

message Book {}
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent

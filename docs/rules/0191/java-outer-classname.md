---
rule:
  aip: 191
  name: [core, '0191', java-outer-classname]
  summary: All proto files must set `option java_outer_classname`.
permalink: /191/java-outer-classname
redirect_from:
  - /0191/java-outer-classname
---

# Java package annotation

This rule enforces that every proto file for a public API surface sets
`option java_outer_classname`, as mandated in [AIP-191][].

## Details

This rule looks at each proto file, and complains if the `java_outer_classname`
file annotation is not present, or set to something other than the common class
name based on the proto's filename.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

option java_package = "com.google.example.v1";
option java_multiple_files = true;
// Needs `option java_outer_classname = "LibraryProto";` or similar.

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
// (-- api-linter: core::0191::java-outer-classname=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

package google.example.v1;

option java_package = "com.google.example.v1";
option java_multiple_files = true;

message Book {}
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent

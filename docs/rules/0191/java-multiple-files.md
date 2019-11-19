---
rule:
  aip: 191
  name: [core, '0191', java-multiple-files]
  summary: All proto files must set `option java_multiple_files = true`.
permalink: /191/java-multiple-files
redirect_from:
  - /0191/java-multiple-files
---

# Java multiple files annotation

This rule enforces that every proto file for a public API surface sets
`option java_multiple_files = true;`, as mandated in [AIP-191][].

## Details

This rule looks at each proto file, and complains if the `java_multiple_files`
file annotation is not present, or set to `false`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto3";

// Needs `option java_multiple_files = true;`.

message Book {}
```

**Correct** code for this rule:

```proto
// Correct.
syntax = "proto3";

option java_multiple_files = true;

message Book {}
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::java-multiple-files=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

message Book {}
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent

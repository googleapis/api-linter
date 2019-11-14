---
rule:
  aip: 191
  name: [core, '0191', proto-version]
  summary: All proto files must use proto3.
permalink: /191/proto-version
redirect_from:
  - /0191/proto-version
---

# Proto3 syntax

This rule enforces that every proto file for a public API surface uses proto3,
as mandated in [AIP-191][].

## Details

This rule looks at each proto file, and complains if the syntax is set to
`proto2` (or missing, which means it defaults to `proto2`).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto2";  // Should be proto3.
```

**Correct** code for this rule:

```proto
// Correct.
syntax = "proto3";
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::proto-version=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto2";
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent

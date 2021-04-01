---
rule:
  aip: 141
  name: [core, '0141', forbidden-types]
  summary: Fields should avoid unsigned integer types.
permalink: /141/forbidden-types
redirect_from:
  - /0141/forbidden-types
---

# Forbidden types

This rule enforces that fields do not use unsigned integer types (because many
programming languages and systems do not support them well), as mandated in
[AIP-141][].

## Details

This rule scans all fields and complains if it sees any of the following types:

- `fixed32`
- `fixed64`
- `uint32`
- `uint64`

It suggests use of `int32` or `int64` instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  uint32 page_count = 2;  // Should be `int32`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  int32 page_count = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0141::forbidden-types=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  uint32 page_count = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-141]: https://aip.dev/141
[aip.dev/not-precedent]: https://aip.dev/not-precedent

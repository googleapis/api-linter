---
rule:
  aip: 146
  name: [core, '0146', any]
  summary: Avoid `google.protobuf.Any` fields.
permalink: /146/standardized-codes
redirect_from:
  - /0146/standardized-codes
---

# Any

This rule discourages the use of `google.protobuf.Any`, as described in
[AIP-146][].

## Details

This rule complains if it sees a `google.protobuf.Any` field. Common packages
(such as `google.api` or `google.longrunning`) are excluded.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // google.protobuf.Any is discouraged.
  google.protobuf.Any contents = 1;
}
```

**Correct** code for this rule:

The correct code is likely to vary substantially by use case. See [AIP-146][]
for details and tradeoffs of various approaches for generic fields.

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0146::any=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  google.protobuf.Any contents = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-146]: https://aip.dev/146
[aip.dev/not-precedent]: https://aip.dev/not-precedent

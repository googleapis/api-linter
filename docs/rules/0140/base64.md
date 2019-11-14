---
rule:
  aip: 140
  name: [core, '0140', base64]
  summary: Base64 fields should use the `bytes` type.
permalink: /140/base64
redirect_from:
  - /0140/base64
---

# Base64 fields

This rule tries to enforce that base64 fields use the `bytes` type, as mandated
by [AIP-140][].

## Details

This rule checks the comments over each field, and if "base64" is mentioned, and yet
the field is a `string`, it suggests using `bytes` instead.

## Examples

### Single word method

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // The base64-encoded checksum.
  string checksum = 2;  // Should be bytes.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // The base64-encoded checksum.
  bytes checksum = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0140::base64=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;

  // The base64-encoded checksum
  string checksum = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-140]: https://aip.dev/140
[aip.dev/not-precedent]: https://aip.dev/not-precedent

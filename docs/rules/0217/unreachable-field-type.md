---
rule:
  aip: 217
  name: [core, '0217', unreachable-field-type]
  summary: The unreachable field should be a repeated string.
permalink: /217/unreachable-field-type
redirect_from:
  - /0217/unreachable-field-type
---

# States

This rule enforces that the `unreachable` field is a `repeated string`, as
mandated in [AIP-217][].

## Details

This rule looks at fields named `unreachable`, and complains if they are
anything other than a `repeated string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  string unreachable_locations = 3;  // Should be repeated.
}
```

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  repeated ErrorInfo unreachable = 3;  // Should be a string.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  repeated string unreachable = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  // (-- api-linter: core::0217::unreachable-field-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated string unreachable_locations = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-217]: https://aip.dev/217
[aip.dev/not-precedent]: https://aip.dev/not-precedent

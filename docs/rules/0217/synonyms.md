---
rule:
  aip: 217
  name: [core, '0217', synonyms]
  summary: The field for unreachable resources should be named `unreachable`.
permalink: /217/synonyms
redirect_from:
  - /0217/synonyms
---

# States

This rule enforces that the field describing unreachable resources in response
messages is named `unreachable` rather than `unreachable_locations`, as
mandated in [AIP-217][].

## Details

This rule iterates over paginated response messages and looks for a field named
`unreachable_locations`, and complains if it finds one, suggesting the use of
`unreachable` instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  repeated string unreachable_locations = 3;
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
  // (-- api-linter: core::0217::synonyms=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated string unreachable_locations = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-217]: https://aip.dev/217
[aip.dev/not-precedent]: https://aip.dev/not-precedent

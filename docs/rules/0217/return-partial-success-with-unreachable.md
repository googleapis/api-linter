---
rule:
  aip: 217
  name: [core, '0217', return-partial-success-with-unreachable]
  summary: The return_partial_success field should be paired with unreachable.
permalink: /217/return-partial-success-with-unreachable
redirect_from:
  - /0217/return-partial-success-with-unreachable
---

# States

This rule enforces that the `return_partial_success` field is paired with a
corresponding `repeated string unreachable` field, as mandated in [AIP-217][].

## Details

This rule looks at methods that have a request field named
`return_partial_success`, and complains if the response does not have a
`repeated string unreachable` field.

## Examples

**Incorrect** code for this rule:

```proto
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  // Incorrect. Missing unreachable response field.
  string return_partial_success = 4;
}

message ListBooksResponse {
    repeated Book books = 1;

    string next_page_token = 2;

    // Incorrect. Missing unreachable field.
}
```

**Correct** code for this rule:

```proto
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  // Correct.
  bool return_partial_success = 4;
}

message ListBooksResponse {
    repeated Book books = 1;

    string next_page_token = 2;

    // Correct.
    repeated string unreachable = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  // (-- api-linter: core::0217::return-partial-success-with-unreachable=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string return_partial_success = 4;
}

message ListBooksResponse {
    repeated Book books = 1;

    string next_page_token = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-217]: https://aip.dev/217
[aip.dev/not-precedent]: https://aip.dev/not-precedent

---
rule:
  aip: 217
  name: [core, '0217', return-partial-success-type]
  summary: The return_partial_success field should be a bool.
permalink: /217/return-partial-success-type
redirect_from:
  - /0217/return-partial-success-type
---

# States

This rule enforces that the `return_partial_success` field is a `bool`, as
mandated in [AIP-217][].

## Details

This rule looks at fields named `return_partial_success`, and complains if they
are anything other than a `bool`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  string return_partial_success = 4;  // Should be bool.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  bool return_partial_success = 4;
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
  // (-- api-linter: core::0217::return-partial-success-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string return_partial_success = 4;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-217]: https://aip.dev/217
[aip.dev/not-precedent]: https://aip.dev/not-precedent

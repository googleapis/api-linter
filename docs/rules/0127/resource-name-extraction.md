---
rule:
  aip: 127
  name: [core, '0127', resource-name-extraction]
  summary: HTTP annotations should extract full resource names into variables.
permalink: /127/resource-name-extraction
redirect_from:
  - /0127/resource-name-extraction
---

# HTTP URI case

This rule enforces that HTTP annotations pull whole resource names into
variables, and not just the ID components, as mandated in [AIP-127][].

## Details

This rule scans all methods and complains if it finds a URI with a variable
whose value is `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  // Should be /v1/{name=publishers/*/books/*}
  get: "/v1/publishers/{publisher_id}/books/{book_id}"
}
```

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  // Should be /v1/{name=publishers/*/books/*}
  get: "/v1/publishers/{publisher_id=*}/books/{book_id=*}"
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::resource-name-extraction=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
  get: "/v1/publishers/{publisher_id}/books/{book_id}"
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-127]: https://aip.dev/127
[aip.dev/not-precedent]: https://aip.dev/not-precedent

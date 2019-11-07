---
rule:
  aip: 234
  name: [core, '0234', http-uri-suffix]
  summary: Batch Update methods must have the correct URI suffix
---

# Batch Update methods: URI suffix

This rule enforces that ` BatchUpdate` methods include the `:batchUpdate` suffix
in the REST URI, as mandated in [AIP-234][].

## Details

This rule looks at any method whose name starts with `BatchUpdate`, and
complains if the HTTP URI does not end with `:batchUpdate`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchUpdate`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0234::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchUpdate`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent

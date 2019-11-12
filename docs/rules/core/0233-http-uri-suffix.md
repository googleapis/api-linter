---
rule:
  aip: 233
  name: [core, '0233', http-uri-suffix]
  summary: Batch Create methods must have the correct URI suffix
---

# Batch Create methods: URI suffix

This rule enforces that ` BatchCreate` methods include the `:batchCreate` suffix
in the REST URI, as mandated in [AIP-233][].

## Details

This rule looks at any method whose name starts with `BatchCreate`, and
complains if the HTTP URI does not end with `:batchCreate`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchCreate`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0233::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchCreate`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent

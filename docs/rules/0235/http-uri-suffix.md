---
rule:
  aip: 235
  name: [core, '0235', http-uri-suffix]
  summary: Batch Delete methods must have the correct URI suffix
permalink: /235/http-uri-suffix
redirect_from:
  - /0235/http-uri-suffix
---

# Batch Delete methods: URI suffix

This rule enforces that `BatchDelete` methods include the `:batchDelete` suffix
in the REST URI, as mandated in [AIP-235][].

## Details

This rule looks at any method whose name starts with `BatchDelete`, and
complains if the HTTP URI does not end with `:batchDelete`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchDelete`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0235::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batch" // Should end with `:batchDelete`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent

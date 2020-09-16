---
rule:
  aip: 235
  name: [core, '0235', http-method]
  summary: Batch Delete methods must use the POST HTTP verb.
permalink: /235/http-method
redirect_from:
  - /0235/http-method
---

# Batch Delete methods: POST HTTP verb

This rule enforces that all `BatchDelete` RPCs use the `POST` HTTP verb, as
mandated in [AIP-235][].

## Details

This rule looks at any RPCs with the name beginning with `BatchDelete`, and
complains if the HTTP verb is anything other than `POST`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{parent=publishers/*}/books:batchDelete" // Should be `post:`.
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
// (-- api-linter: core::0235::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{parent=publishers/*}/books:batchDelete" // Should be `post:`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent

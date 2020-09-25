---
rule:
  aip: 235
  name: [core, '0235', http-body]
  summary: Batch Delete methods should use `*` as the HTTP body.
permalink: /235/http-body
redirect_from:
  - /0235/http-body
---

# Batch Delete methods: HTTP body

This rule enforces that all `BatchDelete` RPCs use `*` as the HTTP `body`, as
mandated in [AIP-235][].

## Details

This rule looks at any RPC methods beginning with `BatchDelete`, and
complains if the HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: ""  // The http body should be "*".
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
// (-- api-linter: core::0235::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent

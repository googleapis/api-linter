---
rule:
  aip: 235
  name: [core, '0235', request-message-name]
  summary: Batch Delete methods must have standardized request message names.
permalink: /235/request-message-name
redirect_from:
  - /0235/request-message-name
---

# Batch Delete methods: Request message

This rule enforces that all `BatchDelete` RPCs have a request message name of
`BatchDelete*Request`, as mandated in [AIP-235][].

## Details

This rule looks at any message beginning with `BatchDelete`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `BatchDeleteBooksRequest`.
rpc BatchDeleteBooks(Book) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:batchDelete"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:batchDelete"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0235::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(Book) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:batchDelete"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent

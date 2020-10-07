---
rule:
  aip: 235
  name: [core, '0235', response-message-name]
  summary: Batch Delete methods must have standardized response message names.
permalink: /235/response-message-name
redirect_from:
  - /0235/response-message-name
---

# Batch Delete methods: Response message

This rule enforces that all `BatchDelete` RPCs have a response message name of
`google.protobuf.Empty` or `BatchDelete*Response`, as mandated in [AIP-235][].

## Details

This rule looks at any RPCs whose name starts with `BatchDelete`, and
complains if the name of the corresponding returned message does not match
`google.protobuf.Empty` or the name of the RPC with the suffix `Response`
appended, the latter of which is mandated in [AIP-235][] for soft-delete
operations.

It also permits a response of `google.longrunning.Operation`; in this case, it
checks the `response_type` in the `google.longrunning.operation_info`
annotation and ensures that _it_ is `google.protobuf.Empty` or corresponds to
the name of the RPC with the suffix `Response` appended, as mandated in
[AIP-151][].

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `google.protobuf.Empty` or `BatchDeleteBooksResponse`.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
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
// (-- api-linter: core::0235::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-151]: https://aip.dev/151
[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent

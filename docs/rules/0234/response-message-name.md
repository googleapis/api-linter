---
rule:
  aip: 234
  name: [core, '0234', response-message-name]
  summary: Batch Update methods must have standardized response message names.
permalink: /234/response-message-name
redirect_from:
  - /0234/response-message-name
---

# Batch Update methods: Response message

This rule enforces that all `BatchUpdate` RPCs have a response message name of
`BatchUpdate*Response`, as mandated in [AIP-234][].

## Details

This rule looks at any RPCs whose name starts with `BatchUpdate`, and
complains if the name of the corresponding returned message does not match the
name of the RPC with the suffix `Response` appended.

It also permits a response of `google.longrunning.Operation`; in this case, it
checks the `response_type` in the `google.longrunning.operation_info`
annotation and ensures that _it_ corresponds to the name of the RPC with the
suffix `Response` appended, as mandated in [AIP-151][].

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `BatchUpdateBooksResponse`.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
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
// (-- api-linter: core::0234::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip-151]: https://aip.dev/151
[aip.dev/not-precedent]: https://aip.dev/not-precedent

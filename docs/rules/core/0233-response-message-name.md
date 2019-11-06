---
rule:
  aip: 233
  name: [core, '0233', response-message-name]
  summary: Batch Create methods must have standardized response message names.
---

# Batch Create methods: Response message

This rule enforces that all `BatchCreate` RPCs have a response message name of
`BatchCreate*Response`, as mandated in [AIP-233][].

## Details

This rule looks at any RPCs whose name starts with `BatchCreate`, and
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
// Should be `BatchCreateBooksResponse`.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
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
// (-- api-linter: core::0233::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip-151]: https://aip.dev/151
[aip.dev/not-precedent]: https://aip.dev/not-precedent

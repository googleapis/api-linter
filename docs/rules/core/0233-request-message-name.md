---
rule:
  aip: 233
  name: [core, '0233', request-message-name]
  summary: Batch Create methods must have standardized request message names.
---

# Batch Create methods: Request message

This rule enforces that all `BatchCreate` RPCs have a request message name of
`BatchCreate*Request`, as mandated in [AIP-233][].

## Details

This rule looks at any message matching beginning with `BatchCreate`, and
complains if the name of the corresponding input message does not match the
name of the RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method input should be "BatchCreateBooksRequest"
rpc BatchCreateBooks(CreateBookRequest) returns (BatchCreateBooksResponse) {
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
// (-- api-linter: core::0233::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBooks(CreateBookRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent

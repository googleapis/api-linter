---
rule:
  aip: 233
  name: [core, '0233', http-body]
  summary: Batch Create methods should use `*` as the HTTP body.
---

# Batch Create methods: HTTP body 

This rule enforces that all `BatchCreate` RPCs use `*` as the HTTP `body`, as 
mandated in [AIP-233][].

## Details

This rule looks at any RPC methods matching beginning with `BatchCreate`, and
complains if HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "book"  // The body should be "*".
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
// (-- api-linter: core::0233::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent

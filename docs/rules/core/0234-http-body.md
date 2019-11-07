---
rule:
  aip: 234
  name: [core, '0234', http-body]
  summary: Batch Update methods should use `*` as the HTTP body.
---

# Batch Update methods: HTTP body 

This rule enforces that all `BatchUpdate` RPCs use `*` as the HTTP `body`, as 
mandated in [AIP-234][].

## Details

This rule looks at any RPC methods matching beginning with `BatchUpdate`, and
complains if HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    // Http body is missing.
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
// (-- api-linter: core::0234::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    // Http body is missing.
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent

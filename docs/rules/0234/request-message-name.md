---
rule:
  aip: 234
  name: [core, '0234', request-message-name]
  summary: Batch Update methods must have standardized request message names.
permalink: /234/request-message-name
redirect_from:
  - /0234/request-message-name
---

# Batch Update methods: Request message

This rule enforces that all `BatchUpdate` RPCs have a request message name of
`BatchUpdate*Request`, as mandated in [AIP-234][].

## Details

This rule looks at any RPCs whose name starts with `BatchUpdate`, and
complains if the name of the corresponding input message does not match the
name of the RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method input should be "BatchUpdateBooksRequest"
rpc BatchUpdateBooks(UpdateBookRequest) returns (BatchUpdateBooksResponse) {
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
// (-- api-linter: core::0234::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBooks(UpdateBookRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent

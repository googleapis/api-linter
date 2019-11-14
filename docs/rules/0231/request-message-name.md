---
rule:
  aip: 231
  name: [core, '0231', request-message-name]
  summary: Batch Get methods must have standardized request message names.
permalink: /231/request-message-name
redirect_from:
  - /0231/request-message-name
---

# Batch Get methods: Request message

This rule enforces that all `BatchGet` RPCs have a request message name of
`BatchGet*Request`, as mandated in [AIP-231][].

## Details

This rule looks at any message matching beginning with `BatchGet`, and
complains if the name of the corresponding input message does not match the
name of the RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method input should be "BatchGetBooksRequest"
rpc BatchGetBooks(GetBookRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchGetBooks(GetBookRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent

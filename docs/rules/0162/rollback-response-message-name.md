---
rule:
  aip: 162
  name: [core, '0162', rollback-response-message-name]
  summary: Rollback methods must have standardized response message names.
permalink: /162/rollback-response-message-name
redirect_from:
  - /0162/rollback-response-message-name
---

# Rollback methods: Response message

This rule enforces that all `Rollback` RPCs have a response message of the
resource, as mandated in [AIP-162][].

## Details

This rule looks at any method beginning with `Rollback`, and complains
if the name of the corresponding output message does not match the name of the
RPC with the prefix `Rollback` removed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should return `Book`.
rpc RollbackBook(RollbackBookRequest) returns (RollbackBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:rollback"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc RollbackBook(RollbackBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:rollback"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::rollback-response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc RollbackBook(RollbackBookRequest) returns (RollbackBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:rollback"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent

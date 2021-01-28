---
rule:
  aip: 162
  name: [core, '0162', commit-request-message-name]
  summary: Commit methods must have standardized request message names.
permalink: /162/commit-request-message-name
redirect_from:
  - /0162/commit-request-message-name
---

# Commit methods: Request message

This rule enforces that all `Commit` RPCs have a request message name of
`Commit*Request`, as mandated in [AIP-162][].

## Details

This rule looks at any method beginning with `Commit`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `CommitBookRequest`.
rpc CommitBook(SaveBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::commit-request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CommitBook(SaveBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent

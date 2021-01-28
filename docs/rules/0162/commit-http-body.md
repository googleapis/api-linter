---
rule:
  aip: 162
  name: [core, '0162', commit-http-body]
  summary: Commit methods should use `*` as the HTTP body.
permalink: /162/commit-http-body
redirect_from:
  - /0162/commit-http-body
---

# Commit methods: HTTP body

This rule enforces that all `Commit` RPCs use `*` as the HTTP `body`, as mandated in
[AIP-162][].

## Details

This rule looks at any method beginning with `Commit`, and complains
if the HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
    // body: "*" should be set.
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
// (-- api-linter: core::0162::commit-http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:commit"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent

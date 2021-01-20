---
rule:
  aip: 162
  name: [core, '0162', commit-http-uri-suffix]
  summary: Commit methods must have the correct URI suffix
permalink: /162/commit-http-uri-suffix
redirect_from:
  - /0162/commit-http-uri-suffix
---

# Commit methods: URI suffix

This rule enforces that `Commit` methods include the `:commit` suffix
in the REST URI, as mandated in [AIP-162][].

## Details

This rule looks at any method beginning with `Commit`, and
complains if the HTTP URI does not end with `:commit`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CommitBook(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:save"  // Should end with `:commit`
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
// (-- api-linter: core::0162::commit-http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc Commit(CommitBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:save"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent

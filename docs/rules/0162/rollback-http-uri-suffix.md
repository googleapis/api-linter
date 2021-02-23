---
rule:
  aip: 162
  name: [core, '0162', rollback-http-uri-suffix]
  summary: Rollback methods must have the correct URI suffix
permalink: /162/rollback-http-uri-suffix
redirect_from:
  - /0162/rollback-http-uri-suffix
---

# Rollback methods: URI suffix

This rule enforces that `Rollback` methods include the `:rollback` suffix
in the REST URI, as mandated in [AIP-162][].

## Details

This rule looks at any method beginning with `Rollback`, and
complains if the HTTP URI does not end with `:rollback`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc RollbackBook(RollbackBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undo"  // Should end with `:rollback`
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
// (-- api-linter: core::0162::rollback-http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc RollbackBook(RollbackBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undo"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent

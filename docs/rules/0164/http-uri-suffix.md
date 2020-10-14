---
rule:
  aip: 164
  name: [core, '0164', http-uri-suffix]
  summary: Undelete methods must have the correct URI suffix
permalink: /164/http-uri-suffix
redirect_from:
  - /0164/http-uri-suffix
---

# Undelete methods: URI suffix

This rule enforces that `Undelete` methods include the `:undelete` suffix
in the REST URI, as mandated in [AIP-164][].

## Details

This rule looks at any method whose name starts with `Undelete`, and
complains if the HTTP URI does not end with `:undelete`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:restore" // Should end with `:undelete`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0164::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:restore" // Should end with `:undelete`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent

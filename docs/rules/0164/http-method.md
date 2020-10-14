---
rule:
  aip: 164
  name: [core, '0164', http-method]
  summary: Undelete methods must use the POST HTTP verb.
permalink: /164/http-method
redirect_from:
  - /0164/http-method
---

# Undelete methods: POST HTTP verb

This rule enforces that all `Undelete` RPCs use the `POST` HTTP verb, as
mandated in [AIP-164][].

## Details

This rule looks at any message beginning with `Undelete`, and complains
if the HTTP verb is anything other than `POST`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:undelete"  // Should be `post:`.
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
// (-- api-linter: core::0164::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:undelete"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent

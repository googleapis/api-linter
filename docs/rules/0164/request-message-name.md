---
rule:
  aip: 164
  name: [core, '0164', request-message-name]
  summary: Undelete methods must have standardized request message names.
permalink: /164/request-message-name
redirect_from:
  - /0164/request-message-name
---

# Undelete methods: Request message

This rule enforces that all `Undelete` RPCs have a request message name of
`Undelete*Request`, as mandated in [AIP-164][].

## Details

This rule looks at any message beginning with `Undelete`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `UndeleteBookRequest`.
rpc UndeleteBook(Book) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
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
// (-- api-linter: core::0164::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(Book) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent

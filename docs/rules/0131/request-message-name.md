---
rule:
  aip: 131
  name: [core, '0131', request-message-name]
  summary: Get methods must have standardized request message names.
permalink: /131/request-message-name
redirect_from:
  - /0131/request-message-name
---

# Get methods: Request message

This rule enforces that all `Get` RPCs have a request message name of
`Get*Request`, as mandated in [AIP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the name of the corresponding input message does not match the name of the RPC
with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookReq) returns (Book) {  // Should be `GetBookRequest`.
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookReq) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent

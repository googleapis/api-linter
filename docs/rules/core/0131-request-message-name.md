---
rule:
  aip: 126
  name: [core, '0126', upper-snake-values]
  summary: All enum values must be in upper snake case.
---

# Upper snake case values

This rule enforces that all `Get*` RPCs have a message name of `Get*Request`,
as mandated in [AIP-131][].

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
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::request-message-name=disabled
//     aip.dev/not-precedent: Named this way for historical reasons. --)
rpc GetBook(GetBookReq) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent

---
rule:
  aip: 136
  name: [core, '0136', http-body]
  summary: Custom methods must have the HTTP body set to `*`.
---

# Custom methods: HTTP body

This rule enforces that all custom methods set the HTTP `body` to `*`, as
mandated in [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and complains 
if the HTTP `body` field is not set to the resource being created.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:checkout"
    body: "book"  // This should be "*".
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:checkout"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:checkout"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent

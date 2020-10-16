---
rule:
  aip: 136
  name: [core, '0136', standard-methods-only]
  summary: Declarative-friendly resources should eschew custom methods.
permalink: /136/standard-methods-only
redirect_from:
  - /0136/standard-methods-only
---

# Declarative: Standard methods only

This rule enforces that declarative-friendly resources do not use custom
methods, as discussed in [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and if the
associated resource is a declarative-friendly resource, complains about the use
of a custom method.

## Examples

### Verb only

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that book is declarative-friendly...
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

**Correct** code for this rule:

**Important:** In general, declarative-friendly resources **should not** use
custom methods; however, some declarative-friendly resources **may** have
one-off, truly imperative methods that do not expect support in imperative
tooling. This can be indicated to the linter using the comment shown below.

```proto
// Correct.
// (-- Imperative only. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::standard-methods-only=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

**Note:** For truly imperative-only methods, you can use the comment form shown
above.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[http-name-variable]: ./http-name-variable.md
[http-parent-variable]: ./http-parent-variable.md
